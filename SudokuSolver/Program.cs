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
            SudokuSolver.Logics.Solver.Call();

            //int currentInputIndex = 0;
            //Random r = new Random();

            //Console.WriteLine("Loading solver...");
            //Solver s = new Solver();
            //var listOfPer = SolverHelper.Permute(new int[] { 1,2,3,4,5,6,7,8,9 });
            //Console.WriteLine($"Possible line solutions: {listOfPer.Count}");
            //Console.WriteLine($"Possible line length: {listOfPer[0].Count}");
            //Console.WriteLine($"Check count at {1}: {listOfPer.Where(x => x[0] == 1).Count()}");

            //Console.WriteLine("Calculating possible permutations...");
            //Console.WriteLine("------------------------------");
            //for (int l = 0; l < 9; l++)
            //{
            //    SLine sl = new SLine();
            //    sl.data = input.GetNext(l, 0);
            //    sl.possiblePermutations = sl.data.GetPermutations(listOfPer);
            //    sLines.Add(sl);

            //    Console.WriteLine($"Line {l+1} has {sl.possiblePermutations.Count} possible solutions");
            //}
            //Console.WriteLine("------------------------------");

            //for (int l = 0; l < 9; l++)
            //{
            //    SLine sl = new SLine();
            //    sl.data = input.GetNext(l, 1);
            //    sl.possiblePermutations = sl.data.GetPermutations(listOfPer);
            //    sLines.Add(sl);

            //    Console.WriteLine($"Line {l + 1} has {sl.possiblePermutations.Count} possible solutions");
            //}
            //Console.WriteLine("------------------------------");

            //for (int l = 0; l < 9; l++)
            //{
            //    SLine sl = new SLine();
            //    sl.data = input.GetNext(l, 2);
            //    sl.possiblePermutations = sl.data.GetPermutations(listOfPer);
            //    sLines.Add(sl);

            //    Console.WriteLine($"Line {l + 1} has {sl.possiblePermutations.Count} possible solutions");
            //}
            //Console.WriteLine("------------------------------");

            //Console.WriteLine("permutations calculated.");
            //Console.WriteLine("Calculating possible combos...");



            //string solvedStr = s.Solve(inputSolved) ? "yes" : "no";
            //Console.WriteLine($"Solved: {solvedStr}");
            //Console.ReadLine();
        }
    }
}
