#!/usr/bin/env python3

import argparse
import html
import json
import os
import subprocess
import time
from pathlib import Path


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


def render_svg(results: list[dict], commit_sha: str) -> str:
    short_sha = commit_sha[:12] if commit_sha else "local"
    columns = [
        (28, "Difficulty"),
        (210, "Givens"),
        (330, "Completions"),
        (490, "Analysis"),
        (625, "Solve steps"),
        (770, "Full solve"),
        (915, "Result"),
    ]
    row_y = 158
    row_height = 34
    height = row_y + row_height * len(results) + 30

    parts = [
        f'<svg xmlns="http://www.w3.org/2000/svg" width="1040" height="{height}" viewBox="0 0 1040 {height}" role="img" aria-label="Latest SudokuSolver benchmark results">',
        "<style>",
        ".bg{fill:#ffffff}.title{font:700 22px -apple-system,BlinkMacSystemFont,Segoe UI,sans-serif;fill:#1f2328}",
        ".meta{font:13px -apple-system,BlinkMacSystemFont,Segoe UI,sans-serif;fill:#57606a}",
        ".head{font:600 13px -apple-system,BlinkMacSystemFont,Segoe UI,sans-serif;fill:#1f2328}",
        ".cell{font:13px ui-monospace,SFMono-Regular,Consolas,monospace;fill:#1f2328}",
        ".line{stroke:#d0d7de;stroke-width:1}",
        "</style>",
        f'<rect class="bg" width="1040" height="{height}" rx="8"/>',
        '<text class="title" x="28" y="36">SudokuSolver automatic benchmarks</text>',
        '<text class="meta" x="28" y="62">Exact analysis + probability-guided full solve on a GitHub-hosted Ubuntu runner</text>',
        f'<text class="meta" x="28" y="84">Source commit: {html.escape(short_sha)} · timings vary between runners</text>',
        '<line class="line" x1="28" y1="106" x2="1012" y2="106"/>',
    ]

    for x, label in columns:
        parts.append(f'<text class="head" x="{x}" y="132">{html.escape(label)}</text>')
    parts.append('<line class="line" x1="28" y1="142" x2="1012" y2="142"/>')

    for index, result in enumerate(results):
        y = row_y + index * row_height
        values = [
            result["difficulty"],
            str(result["givens"]),
            str(result["remaining_solutions"]),
            format_duration(result["analysis_seconds"]),
            str(result["solve_steps"]),
            format_duration(result["solve_seconds"]),
            "Solved",
        ]
        for (x, _), value in zip(columns, values):
            parts.append(f'<text class="cell" x="{x}" y="{y}">{html.escape(value)}</text>')
        parts.append(f'<line class="line" x1="28" y1="{y + 12}" x2="1012" y2="{y + 12}"/>')

    parts.append("</svg>")
    return "\n".join(parts) + "\n"


def main() -> int:
    parser = argparse.ArgumentParser(description="Run SudokuSolver benchmarks and generate README-display results")
    parser.add_argument("--binary", type=Path, required=True)
    parser.add_argument("--fixtures", type=Path, default=Path("benchmarks/puzzles.json"))
    parser.add_argument("--results-json", type=Path, default=Path("benchmark-results.json"))
    parser.add_argument("--results-svg", type=Path, default=Path("benchmark-results.svg"))
    parser.add_argument("--timeout", type=int, default=300, help="timeout in seconds per solver invocation")
    args = parser.parse_args()

    if args.timeout <= 0:
        parser.error("--timeout must be greater than zero")
    if not args.binary.is_file():
        parser.error(f"solver binary not found: {args.binary}")

    fixtures = load_fixtures(args.fixtures)
    results = [benchmark_fixture(args.binary, fixture, args.timeout) for fixture in fixtures]

    args.results_json.write_text(json.dumps(results, indent=2) + "\n", encoding="utf-8")
    args.results_svg.write_text(render_svg(results, os.environ.get("GITHUB_SHA", "")), encoding="utf-8")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
