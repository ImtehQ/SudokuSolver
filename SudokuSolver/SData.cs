using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using SudokuSolver;

namespace SudokuSolver
{
    public class SData
    {
        public int[] data = new int[9];

        public List<List<int>> possiblePermutations;

        public bool isValid()
        {
            if (!data.Contains('0'))
                return (data.Distinct().Count() == 9);
            return false;
        }
    }
}
