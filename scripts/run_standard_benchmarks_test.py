import tempfile
import unittest
from pathlib import Path

from scripts import run_standard_benchmarks as standard


class StandardBenchmarkTests(unittest.TestCase):
    def test_default_sample_size_is_ten(self):
        self.assertEqual(standard.DEFAULT_SAMPLE_SIZE, 10)

    def test_normalize_line_accepts_metadata(self):
        puzzle = "4.....8.5.3..........7......2.....6.....8.4......1.......6.3.7.5..2.....1.4......"
        self.assertEqual(standard.normalize_line(puzzle + " rating=hard"), puzzle)

    def test_write_sample_is_deterministic(self):
        puzzles = [str(index).zfill(81) for index in range(20)]
        with tempfile.TemporaryDirectory() as temp_dir:
            first = Path(temp_dir) / "first.txt"
            second = Path(temp_dir) / "second.txt"
            standard.write_sample(puzzles, 5, 123, first)
            standard.write_sample(puzzles, 5, 123, second)
            self.assertEqual(first.read_text(encoding="utf-8"), second.read_text(encoding="utf-8"))
            self.assertEqual(len(first.read_text(encoding="utf-8").splitlines()), 5)

    def test_format_rate(self):
        self.assertEqual(standard.format_rate(2_500_000), "2.50M/s")
        self.assertEqual(standard.format_rate(12_345), "12.3k/s")
        self.assertEqual(standard.format_rate(321), "321/s")
        self.assertEqual(standard.format_rate(4.25), "4.25/s")

    def test_timeout_result_is_renderable(self):
        timed_out = standard.timeout_result(49158, 300, 300.1)
        self.assertEqual(standard.metric_text(timed_out), ">300s timeout")
        self.assertEqual(standard.speedup_text(timed_out, {"puzzles_per_second": 1000.0}), "—")

    def test_render_svg_contains_rates_timeout_and_commit(self):
        rows = [
            {
                "dataset": "Fixture",
                "our_exact": {"puzzles": 10, "puzzles_per_second": 100.0, "timed_out": False},
                "tdoku_unique": {"puzzles": 10, "puzzles_per_second": 500.0, "timed_out": False},
                "our_probability_solve": {"puzzles": 5, "puzzles_per_second": 20.0, "timed_out": False},
            },
            {
                "dataset": "Full corpus",
                "our_exact": standard.timeout_result(49158, 300, 300.0),
                "tdoku_unique": {"puzzles": 49158, "puzzles_per_second": 1000.0, "timed_out": False},
                "our_probability_solve": None,
            },
        ]
        svg = standard.render_svg(rows, "1234567890abcdef")
        self.assertIn("1234567890ab", svg)
        self.assertIn("5.0x", svg)
        self.assertIn("100/s", svg)
        self.assertIn("500/s", svg)
        self.assertIn("&gt;300s timeout", svg)


if __name__ == "__main__":
    unittest.main()
