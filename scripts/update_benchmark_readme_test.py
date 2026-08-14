import json
import tempfile
import unittest
from pathlib import Path

from scripts import update_benchmark_readme as benchmark


class BenchmarkReadmeTests(unittest.TestCase):
    def test_format_duration(self):
        self.assertEqual(benchmark.format_duration(0.125), "125 ms")
        self.assertEqual(benchmark.format_duration(1.23456), "1.235 s")

    def test_render_svg_contains_results(self):
        svg = benchmark.render_svg(
            [
                {
                    "difficulty": "Impossible",
                    "givens": 21,
                    "remaining_solutions": "1",
                    "analysis_seconds": 0.163,
                    "solve_steps": 60,
                    "solve_seconds": 0.212,
                    "solved": True,
                    "final_grid": "",
                }
            ],
            "abcdef1234567890",
        )
        self.assertIn("SudokuSolver automatic benchmarks", svg)
        self.assertIn("abcdef123456", svg)
        self.assertIn("Impossible", svg)
        self.assertIn("163 ms", svg)
        self.assertIn("212 ms", svg)
        self.assertIn("Solved", svg)

    def test_load_fixtures_requires_expected_profiles(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            path = Path(temp_dir) / "puzzles.json"
            fixtures = [
                {"difficulty": label, "puzzle": "0" * 81}
                for label in ["Easy", "Medium", "Hard", "Impossible"]
            ]
            path.write_text(json.dumps(fixtures), encoding="utf-8")
            loaded = benchmark.load_fixtures(path)
            self.assertEqual([item["difficulty"] for item in loaded], ["Easy", "Medium", "Hard", "Impossible"])

    def test_load_fixtures_rejects_invalid_puzzle(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            path = Path(temp_dir) / "puzzles.json"
            fixtures = [
                {"difficulty": label, "puzzle": "0" * 81}
                for label in ["Easy", "Medium", "Hard", "Impossible"]
            ]
            fixtures[-1]["puzzle"] = "bad"
            path.write_text(json.dumps(fixtures), encoding="utf-8")
            with self.assertRaises(ValueError):
                benchmark.load_fixtures(path)


if __name__ == "__main__":
    unittest.main()
