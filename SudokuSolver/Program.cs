using System;
using System.Collections.Generic;
using System.Linq;

namespace SudokuSolver
{
    class Program
    {
        public static int[] input = new int[] 
        { 
            3,0,1,9,0,8,0,0,5,
            5,0,0,0,0,3,0,0,0,
            0,6,0,0,5,0,0,8,0,
            1,9,6,0,0,0,0,0,0,
            0,2,3,0,0,0,4,7,0,
            0,0,0,0,0,0,2,9,1,
            0,3,0,0,4,0,0,1,0,
            0,0,4,7,0,0,0,0,8,
            6,0,5,3,0,1,9,0,0
        };

        public static int[] inputSolved = new int[]
{
            3,4,1,9,6,8,7,2,5,
            5,8,9,2,7,3,1,6,4,
            7,6,2,1,5,4,3,8,9,
            1,9,6,4,2,7,8,5,3,
            8,2,3,5,1,9,4,7,0,
            4,5,7,8,3,6,2,9,1,
            9,3,8,6,4,2,5,1,7,
            2,1,4,7,9,5,6,3,8,
            6,7,5,3,8,1,9,4,2
};

        //public static List<SLine> sLines = new List<SLine>();

        static void Main(string[] args)
        {
            SudokuSolver.Solver.Call();

        }
    }
}
