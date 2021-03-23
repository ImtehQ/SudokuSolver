using System;
using System.Collections.Generic;
using System.Linq;

namespace SudokuSolver
{
    public enum Dimension
    {
        Line,
        Colum,
        Block
    }
    public enum SolverType
    {
        SolveLogical,
        SolveGuessing,
        SolveEmpty
    }

    public static class Solver
    {
        public static bool showDebug = false;
        public static int[] input = new int[]
{
            0,0,4,3,0,0,2,0,9,
            0,0,5,0,0,9,0,0,1,
            0,7,0,0,6,0,0,4,3,
            0,0,6,0,0,2,0,8,7,
            1,9,0,0,0,7,4,0,0,
            0,5,0,0,8,3,0,0,0,
            6,0,0,0,0,0,1,0,5,
            0,0,3,5,0,8,6,9,0,
            0,4,2,9,1,0,3,0,0
};

        public static int[] inputSolved = new int[]
{
            8,6,4,3,7,1,2,5,9,
            3,2,5,8,4,9,7,6,1,
            9,7,1,2,6,5,8,4,3,
            4,3,6,1,9,2,5,8,7,
            1,9,8,6,5,7,4,3,2,
            2,5,7,4,8,3,9,1,6,
            6,8,9,7,3,4,1,2,5,
            7,1,3,5,2,8,6,9,4,
            5,4,2,9,1,6,3,7,8
};

        public static int[] Hardest = new int[]
        {
            8,0,0,0,0,0,0,0,0,
            0,0,3,6,0,0,0,0,0,
            0,7,0,0,9,0,2,0,0,
            0,5,0,0,0,7,0,0,0,
            0,0,0,0,4,5,7,0,0,
            0,0,0,1,0,0,0,3,0,
            0,0,1,0,0,0,0,6,8,
            0,0,8,5,0,0,0,1,0,
            0,9,0,0,0,0,4,0,0
        };

        static SDataSet sData;
        static int foundTotaal = 0;
        public static void Call()
        {
            input = Hardest;
            if (sData == null) sData = new SDataSet();
            sData.Reset();
            Random r = new Random();

            Console.WriteLine("resetting solver...");

            sData.Permutations = SolverHelper.Permute(new int[] { 1, 2, 3, 4, 5, 6, 7, 8, 9 });

            if (showDebug)
                Console.WriteLine($"Possible line solutions: {sData.Permutations.Count}");
            if (showDebug)
                Console.WriteLine($"Possible line length: {sData.Permutations[0].Count}");
            if (showDebug)
                Console.WriteLine($"Check count at {1}: {sData.Permutations.Where(x => x[0] == 1).Count()}");

            Console.WriteLine("Reading sudoku input...");
            sData.ReadSudokuData(input, Dimension.Line);
            sData.ReadSudokuData(input, Dimension.Colum);
            sData.ReadSudokuData(input, Dimension.Block);

            if (showDebug)
                Console.WriteLine("Converted to Sudoku data set model.");
            if (showDebug)
                Console.WriteLine("------------------------------");
            if (showDebug)
                Console.WriteLine("Calculating possible permutations...");

            sData.GetPermutations();

            if (showDebug)
                Console.WriteLine("Permutations calculated.");
            if (showDebug)
                Console.WriteLine("------------------------------");


            int[] li;
            int[] ci;
            int[] bi;




            for (int o = 0; o < 10; o++)
            {
                for (int i = 0; i < 81; i++)
                {
                    sData.CheckISValid(i);
                }
            }



            for (int i = 0; i < 9; i++)
            {
                Console.WriteLine($"==============================================");
                Console.WriteLine($"Y: {i} ==" +
                                    $" li: {sData.sdLines[i].possiblePermutations.Count()}," +
                                    $" ci: {sData.sdCols[i].possiblePermutations.Count()}," +
                                    $" bi: {sData.sdBlocks[i].possiblePermutations.Count()}");

                for (int u = 1; u < 10; u++)
                {
                    Console.WriteLine($"Checking {u}, Count: {sData.sdBlocks[i].possiblePermutations.Count(x => x[0] == u)}");
                }
                
            }

            

            Console.ReadLine();
        }
    }
}