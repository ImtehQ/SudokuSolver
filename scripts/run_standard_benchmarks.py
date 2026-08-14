#!/usr/bin/env python3

import argparse
import html
import json
import os
import random
import subprocess
import tempfile
import time
from pathlib import Path

T_DOKU_COMMIT = "af426180dc53aef89b82868e7b3fdfcf42165654"
SAMPLE_SEED = 20260814


def normalize_line(line: str) -> str | None:
    line = line.strip()
    if not line or line.startswith("#"):
        return None
    token = line.split()[0]
    if len(token) < 81:
        raise ValueError(f"puzzle token has {len(token)} cells, want at least 81")
    token = token[:81]
    if any(char not in "1234567890." for char in token):
        raise ValueError("puzzle contains an unsupported character")
    return token


def load_puzzles(path: Path) -> list[str]:
    puzzles = []
    for line_number, line in enumerate(path.read_text(encoding="utf-8").splitlines(), start=1):
        try:
            puzzle = normalize_line(line)
        except ValueError as exc:
            raise ValueError(f"{path}:{line_number}: {exc}") from exc
        if puzzle is not None:
            puzzles.append(puzzle)
    if not puzzles:
        raise ValueError(f"{path} contains no puzzles")
    return puzzles


def write_sample(puzzles: list[str], size: int, seed: int, path: Path) -> Path:
    if size <= 0:
        raise ValueError("sample size must be positive")
    if size > len(puzzles):
        size = len(puzzles)
    if size == len(puzzles):
        selected = puzzles
    else:
        rng = random.Random(seed)
        indices = sorted(rng.sample(range(len(puzzles)), size))
        selected = [puzzles[index] for index in indices]
    path.write_text("\n".join(selected) + "\n", encoding="utf-8")
    return path


def run_json(command: list[str], timeout: int, stdin_path: Path | None = None) -> tuple[dict, float]:
    started = time.perf_counter()
    if stdin_path is None:
        completed = subprocess.run(command, capture_output=True, text=True, timeout=timeout, check=False)
    else:
        with stdin_path.open("r", encoding="utf-8") as input_file:
            completed = subprocess.run(
                command,
                stdin=input_file,
                capture_output=True,
                text=True,
                timeout=timeout,
                check=False,
            )
    elapsed = time.perf_counter() - started
    if completed.returncode != 0:
        raise RuntimeError(
            f"command failed ({completed.returncode}): {' '.join(command)}\n"
            f"stdout:\n{completed.stdout}\nstderr:\n{completed.stderr}"
        )
    try:
        return json.loads(completed.stdout), elapsed
    except json.JSONDecodeError as exc:
        raise RuntimeError(f"command did not return JSON: {' '.join(command)}\n{completed.stdout}") from exc


def run_ours(binary: Path, dataset: str, path: Path, mode: str, timeout: int) -> dict:
    result, wall_seconds = run_json(
        [
            str(binary),
            "--input",
            str(path),
            "--dataset",
            dataset,
            "--mode",
            mode,
            "--require-unique=true",
        ],
        timeout,
    )
    result["wall_seconds"] = wall_seconds
    return result


def run_tdoku(binary: Path, path: Path, timeout: int) -> dict:
    result, wall_seconds = run_json([str(binary)], timeout, stdin_path=path)
    puzzles = int(result["puzzles"])
    if puzzles <= 0 or int(result["unique"]) != puzzles:
        raise RuntimeError(f"Tdoku uniqueness adapter returned invalid result: {result}")
    result["wall_seconds"] = wall_seconds
    result["puzzles_per_second"] = puzzles / wall_seconds if wall_seconds else 0.0
    return result


def format_rate(value: float) -> str:
    if value >= 1_000_000:
        return f"{value / 1_000_000:.2f}M/s"
    if value >= 1_000:
        return f"{value / 1_000:.1f}k/s"
    if value >= 100:
        return f"{value:.0f}/s"
    if value >= 10:
        return f"{value:.1f}/s"
    return f"{value:.2f}/s"


def render_svg(results: list[dict], commit_sha: str) -> str:
    short_sha = commit_sha[:12] if commit_sha else "local"
    width = 1220
    row_y = 180
    row_height = 42
    height = row_y + row_height * len(results) + 48
    columns = [
        (28, "Dataset"),
        (330, "Count N"),
        (425, "SudokuSolver exact"),
        (610, "Tdoku unique"),
        (775, "Tdoku speedup"),
        (930, "Probability solve"),
        (1110, "Solve N"),
    ]

    parts = [
        f'<svg xmlns="http://www.w3.org/2000/svg" width="{width}" height="{height}" viewBox="0 0 {width} {height}" role="img" aria-label="SudokuSolver standard benchmark comparison">',
        "<style>",
        ".bg{fill:#ffffff}.title{font:700 22px -apple-system,BlinkMacSystemFont,Segoe UI,sans-serif;fill:#1f2328}",
        ".meta{font:13px -apple-system,BlinkMacSystemFont,Segoe UI,sans-serif;fill:#57606a}",
        ".head{font:600 13px -apple-system,BlinkMacSystemFont,Segoe UI,sans-serif;fill:#1f2328}",
        ".cell{font:13px ui-monospace,SFMono-Regular,Consolas,monospace;fill:#1f2328}",
        ".line{stroke:#d0d7de;stroke-width:1}",
        "</style>",
        f'<rect class="bg" width="{width}" height="{height}" rx="8"/>',
        '<text class="title" x="28" y="36">Standard Sudoku benchmark comparison</text>',
        '<text class="meta" x="28" y="62">Same GitHub-hosted Ubuntu runner · unique-puzzle proof: exact count vs pinned Tdoku limit=2</text>',
        f'<text class="meta" x="28" y="84">SudokuSolver source: {html.escape(short_sha)} · Tdoku: {T_DOKU_COMMIT[:12]} · deterministic samples: seed {SAMPLE_SEED}</text>',
        '<text class="meta" x="28" y="106">Probability solve is extra work unique to this project: exact distribution + highest-probability choice repeated to completion.</text>',
        '<line class="line" x1="28" y1="128" x2="1192" y2="128"/>',
    ]
    for x, label in columns:
        parts.append(f'<text class="head" x="{x}" y="154">{html.escape(label)}</text>')
    parts.append('<line class="line" x1="28" y1="164" x2="1192" y2="164"/>')

    for index, row in enumerate(results):
        y = row_y + index * row_height
        our_rate = float(row["our_exact"]["puzzles_per_second"])
        tdoku_rate = float(row["tdoku_unique"]["puzzles_per_second"])
        speedup = tdoku_rate / our_rate if our_rate else 0.0
        solve_rate = float(row["our_probability_solve"]["puzzles_per_second"])
        values = [
            row["dataset"],
            str(row["our_exact"]["puzzles"]),
            format_rate(our_rate),
            format_rate(tdoku_rate),
            f"{speedup:.1f}x",
            format_rate(solve_rate),
            str(row["our_probability_solve"]["puzzles"]),
        ]
        for (x, _), value in zip(columns, values):
            parts.append(f'<text class="cell" x="{x}" y="{y}">{html.escape(value)}</text>')
        parts.append(f'<line class="line" x1="28" y1="{y + 14}" x2="1192" y2="{y + 14}"/>')

    parts.append("</svg>")
    return "\n".join(parts) + "\n"


def benchmark_dataset(
    label: str,
    full_path: Path,
    solve_path: Path,
    our_binary: Path,
    tdoku_binary: Path,
    timeout: int,
) -> dict:
    our_exact = run_ours(our_binary, label, full_path, "count", timeout)
    tdoku_unique = run_tdoku(tdoku_binary, full_path, timeout)
    if our_exact["puzzles"] != tdoku_unique["puzzles"]:
        raise RuntimeError(f"{label}: SudokuSolver and Tdoku processed different puzzle counts")
    our_probability_solve = run_ours(our_binary, label, solve_path, "solve", timeout)
    return {
        "dataset": label,
        "our_exact": our_exact,
        "tdoku_unique": tdoku_unique,
        "our_probability_solve": our_probability_solve,
    }


def main() -> int:
    parser = argparse.ArgumentParser(description="Benchmark SudokuSolver against standard public Sudoku corpora")
    parser.add_argument("--our-binary", type=Path, required=True)
    parser.add_argument("--tdoku-binary", type=Path, required=True)
    parser.add_argument("--top95", type=Path, required=True)
    parser.add_argument("--min17", type=Path, required=True)
    parser.add_argument("--forum-hardest", type=Path, required=True)
    parser.add_argument("--sample-size", type=int, default=1000)
    parser.add_argument("--timeout", type=int, default=3600)
    parser.add_argument("--results-json", type=Path, default=Path("standard-benchmark-results.json"))
    parser.add_argument("--results-svg", type=Path, default=Path("standard-benchmark-results.svg"))
    args = parser.parse_args()

    if args.sample_size <= 0:
        parser.error("--sample-size must be positive")
    if args.timeout <= 0:
        parser.error("--timeout must be positive")
    for binary in [args.our_binary, args.tdoku_binary]:
        if not binary.is_file():
            parser.error(f"binary not found: {binary}")

    top95 = load_puzzles(args.top95)
    min17 = load_puzzles(args.min17)
    forum = load_puzzles(args.forum_hardest)
    if len(top95) != 95:
        raise RuntimeError(f"Norvig Top95 contains {len(top95)} puzzles, expected 95")
    if len(min17) != 49158:
        raise RuntimeError(f"Tdoku 17-clue corpus contains {len(min17)} puzzles, expected 49158")
    if len(forum) < args.sample_size:
        raise RuntimeError(f"forum-hardest corpus only contains {len(forum)} puzzles")

    with tempfile.TemporaryDirectory(prefix="sudokusolver-standard-") as temp_dir:
        temp = Path(temp_dir)
        top95_path = write_sample(top95, len(top95), SAMPLE_SEED, temp / "top95.txt")
        min17_solve = write_sample(min17, args.sample_size, SAMPLE_SEED, temp / "min17-solve.txt")
        forum_sample = write_sample(forum, args.sample_size, SAMPLE_SEED, temp / "forum-hardest-sample.txt")

        results = [
            benchmark_dataset("Norvig Top95", top95_path, top95_path, args.our_binary, args.tdoku_binary, args.timeout),
            benchmark_dataset("Tdoku 17-clue (all)", args.min17, min17_solve, args.our_binary, args.tdoku_binary, args.timeout),
            benchmark_dataset(
                f"Forum hardest 11+ (sample {args.sample_size})",
                forum_sample,
                forum_sample,
                args.our_binary,
                args.tdoku_binary,
                args.timeout,
            ),
        ]

    payload = {
        "sudokusolver_commit": os.environ.get("GITHUB_SHA", ""),
        "tdoku_commit": T_DOKU_COMMIT,
        "sample_seed": SAMPLE_SEED,
        "sample_size": args.sample_size,
        "datasets": results,
    }
    args.results_json.write_text(json.dumps(payload, indent=2) + "\n", encoding="utf-8")
    args.results_svg.write_text(render_svg(results, payload["sudokusolver_commit"]), encoding="utf-8")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
