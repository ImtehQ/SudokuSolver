#!/usr/bin/env python3

import argparse
import json
import os
import subprocess
import time
from pathlib import Path

START_MARKER = "<!-- benchmark-results:start -->"
END_MARKER = "<!-- benchmark-results:end -->"


def format_duration(seconds: float) -> str:
    if seconds < 1:
        return f"{seconds * 1000:.0f} ms"
    return f"{seconds:.3f} s"


def load_fixtures(path: Path) -> list[dict[str, str]]:
    fixtures = json.loads(path.read_text(encoding="utf-8"))
    if not isinstance(fixtures, list) or len(fixtures) != 4:
        raise ValueError("benchmark fixture file must contain exactly four puzzles")

    expected = ["Easy", "Medium", "Hard", "Impossible"]
    labels = [item.get("difficulty") for item in fixtures]
    if labels != expected:
        raise ValueError(f"benchmark difficulties must be {expected}, got {labels}")

    for item in fixtures:
        puzzle = item.get("puzzle", "")
        if len(puzzle) != 81 or any(char not in "0123456789" for char in puzzle):
            raise ValueError(f"{item.get('difficulty', 'unknown')} puzzle must be 81 digits")
    return fixtures


def run_json(binary: Path, args: list[str], timeout: int) -> tuple[dict, float]:
    started = time.perf_counter()
    completed = subprocess.run(
        [str(binary), *args],
        capture_output=True,
        text=True,
        timeout=timeout,
        check=False,
    )
    elapsed = time.perf_counter() - started
    if completed.returncode != 0:
        raise RuntimeError(
            f"command failed ({completed.returncode}): {binary} {' '.join(args)}\n"
            f"stdout:\n{completed.stdout}\nstderr:\n{completed.stderr}"
        )
    try:
        return json.loads(completed.stdout), elapsed
    except json.JSONDecodeError as exc:
        raise RuntimeError(f"solver did not return JSON: {completed.stdout}") from exc


def benchmark_fixture(binary: Path, fixture: dict[str, str], timeout: int) -> dict:
    puzzle = fixture["puzzle"]
    analysis, analysis_seconds = run_json(binary, ["--json", puzzle], timeout)
    solve, solve_seconds = run_json(binary, ["--solve", "--json", puzzle], timeout)

    if analysis.get("remaining_solutions") != "1" or not analysis.get("unique_solution"):
        raise RuntimeError(
            f"{fixture['difficulty']} benchmark is expected to be uniquely solvable; "
            f"analysis returned {analysis.get('remaining_solutions')} completion(s)"
        )
    if solve.get("initial_solutions") != analysis.get("remaining_solutions"):
        raise RuntimeError(f"{fixture['difficulty']} analysis/solve solution counts disagree")
    if not solve.get("solved"):
        raise RuntimeError(f"{fixture['difficulty']} benchmark was not solved")

    steps = solve.get("steps", [])
    if any(not step.get("cell", {}).get("guaranteed") for step in steps):
        raise RuntimeError(f"{fixture['difficulty']} benchmark used a non-guaranteed step")

    return {
        "difficulty": fixture["difficulty"],
        "givens": sum(char != "0" for char in puzzle),
        "remaining_solutions": analysis["remaining_solutions"],
        "analysis_seconds": analysis_seconds,
        "solve_steps": len(steps),
        "solve_seconds": solve_seconds,
        "solved": True,
        "final_grid": solve.get("final_grid", ""),
    }


def render_results(results: list[dict], commit_sha: str) -> str:
    short_sha = commit_sha[:12] if commit_sha else "local"
    lines = [
        "These results are generated automatically by the `Sudoku Benchmarks` GitHub Action.",
        "Each fixture must have exactly one valid completion, and every solve step must therefore be a guaranteed 100% choice.",
        "The difficulty names are fixed benchmark profiles, not a universal Sudoku rating standard.",
        "Timings come from a GitHub-hosted Ubuntu runner and will vary between runs.",
        "",
        f"Benchmark source commit: `{short_sha}`",
        "",
        "| Difficulty | Givens | Exact completions | Analysis time | Solve steps | Full solve time | Result |",
        "|---|---:|---:|---:|---:|---:|---|",
    ]
    for result in results:
        lines.append(
            "| {difficulty} | {givens} | {remaining_solutions} | {analysis_time} | "
            "{solve_steps} | {solve_time} | ✅ solved |".format(
                difficulty=result["difficulty"],
                givens=result["givens"],
                remaining_solutions=result["remaining_solutions"],
                analysis_time=format_duration(result["analysis_seconds"]),
                solve_steps=result["solve_steps"],
                solve_time=format_duration(result["solve_seconds"]),
            )
        )
    return "\n".join(lines)


def replace_results_section(readme: str, rendered: str) -> str:
    start = readme.find(START_MARKER)
    end = readme.find(END_MARKER)
    if start == -1 or end == -1 or end < start:
        raise ValueError("README benchmark markers are missing or out of order")

    before = readme[: start + len(START_MARKER)]
    after = readme[end:]
    return f"{before}\n{rendered.rstrip()}\n{after}"


def main() -> int:
    parser = argparse.ArgumentParser(description="Run SudokuSolver benchmarks and update README results")
    parser.add_argument("--binary", type=Path, required=True)
    parser.add_argument("--fixtures", type=Path, default=Path("benchmarks/puzzles.json"))
    parser.add_argument("--readme", type=Path, default=Path("README.md"))
    parser.add_argument("--results-json", type=Path, default=Path("benchmark-results.json"))
    parser.add_argument("--timeout", type=int, default=300, help="timeout in seconds per solver invocation")
    args = parser.parse_args()

    if args.timeout <= 0:
        parser.error("--timeout must be greater than zero")
    if not args.binary.is_file():
        parser.error(f"solver binary not found: {args.binary}")

    fixtures = load_fixtures(args.fixtures)
    results = [benchmark_fixture(args.binary, fixture, args.timeout) for fixture in fixtures]

    args.results_json.write_text(json.dumps(results, indent=2) + "\n", encoding="utf-8")
    current_readme = args.readme.read_text(encoding="utf-8")
    rendered = render_results(results, os.environ.get("GITHUB_SHA", ""))
    args.readme.write_text(replace_results_section(current_readme, rendered), encoding="utf-8")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
