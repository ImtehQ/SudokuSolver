import json
import tempfile
import unittest
from pathlib import Path

from scripts import update_benchmark_readme as benchmark


class BenchmarkReadmeTests(unittest.TestCase):
    def test_format_duration(self):
        self.assertEqual(benchmark.format_duration(0.125), "125 ms")
        self.assertEqual(benchmark.format_duration(1.23456), "1.235 s")

    def test_replace_results_section(self):
        original = (
            "before\n"
            f"{benchmark.START_MARKER}\n"
            "old\n"
            f"{benchmark.END_MARKER}\n"
            "after\n"
        )
        updated = benchmark.replace_results_section(original, "new results")
        self.assertIn("new results", updated)
        self.assertNotIn("\nold\n", updated)
        self.assertTrue(updated.startswith("before\n"))
        self.assertTrue(updated.endswith("after\n"))

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
