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

            if(showDebug)
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



            for (int i = 0; i < 81; i++)
            {
                Console.WriteLine($"-{i}-");
                int[] li = i.getIndexOf(Dimension.Line);
                Console.WriteLine("L:" + sData.sdLines[li[0]].possiblePermutations.Count());
                int[] ci = i.getIndexOf(Dimension.Colum);
                Console.WriteLine("C:"+sData.sdCols[ci[0]].possiblePermutations.Count());
                int[] bi = i.getIndexOf(Dimension.Block);
                Console.WriteLine("B:"+sData.sdBlocks[bi[0]].possiblePermutations.Count());
            }
            for (int i = 0; i < 81; i++)
            {
                sData.CheckISValid(i);
            }

            int p = 0;
            for (int y = 0; y < 9; y++)
            {
                int[] li = p.getIndexOf(Dimension.Line);
                int[] ci = p.getIndexOf(Dimension.Colum);
                int[] bi = p.getIndexOf(Dimension.Block);
                Console.WriteLine($"==============================================");
                Console.WriteLine($"Y: {y} == li: {sData.sdLines[li[0]].possiblePermutations.Count()}, ci: {sData.sdCols[ci[0]].possiblePermutations.Count()}, bi: {sData.sdBlocks[bi[0]].possiblePermutations.Count()}");

                for (int x = 0; x < 9; x++)
                {
                    var liresult = sData.sdLines[li[0]].possiblePermutations.Where(g => g[li[1]] == x+1).ToList();
                    var ciresult = sData.sdCols[ci[0]].possiblePermutations.Where(g => g[ci[1]] == x+1).ToList();
                    var biresult = sData.sdBlocks[bi[0]].possiblePermutations.Where(g => g[bi[1]] == x+1).ToList();

                    List<int> liresultd = new List<int>();
                    List<int> ciresultd = new List<int>();
                    List<int> biresultd = new List<int>();

                    for (int i = 0; i < liresult.Count; i++)
                    {
                        liresultd.Add(liresult[i][x]);
                    }
                    for (int i = 0; i < ciresult.Count; i++)
                    {
                        ciresultd.Add(ciresult[i][x]);
                    }
                    for (int i = 0; i < biresult.Count; i++)
                    {
                        biresultd.Add(biresult[i][x]);
                    }
                    liresultd = liresultd.Distinct().ToList();
                    ciresultd = ciresultd.Distinct().ToList();
                    biresultd = biresultd.Distinct().ToList();

                    Console.WriteLine($"---------------------------------");
                    Console.WriteLine($"li[{li[0]},{li[1]}] == {x + 1} ==  {liresult.Count} | {liresultd.Count} =");
                    Console.WriteLine($"|||||||||||||||||||||||||||||||||||");
                    for (int t = 0; t < liresultd.Count; t++)
                    {
                        Console.WriteLine(liresultd[t]);
                    }
                    Console.WriteLine($"|||||||||||||||||||||||||||||||||||");
                    Console.WriteLine($"ci[{ci[0]},{ci[1]}] == {x + 1} ==  {ciresult.Count} | {ciresultd.Count} = ");
                    for (int t = 0; t < ciresultd.Count; t++)
                    {
                        Console.WriteLine(ciresultd[t]);
                    }
                    Console.WriteLine($"|||||||||||||||||||||||||||||||||||");
                    Console.WriteLine($"bi[{bi[0]},{bi[1]}] == {x + 1} ==  {biresult.Count} | {biresultd.Count} = ");
                    for (int t = 0; t < ciresultd.Count; t++)
                    {
                        Console.WriteLine(ciresultd[t]);
                    }
                    Console.WriteLine($"---------------------------------");
                    p++;
                }
            }


            int[] smallestFound = new int[3];
            smallestFound[0] = 50000;

            for (int i = 0; i < 81; i++)
            {
                int[] li = i.getIndexOf(Dimension.Line);
                int[] ci = i.getIndexOf(Dimension.Colum);
                int[] bi = i.getIndexOf(Dimension.Block);

      

                if (sData.sdLines[li[0]].possiblePermutations.Count() < smallestFound[0])
                {

                }

                if (sData.sdCols[ci[0]].possiblePermutations.Count() < smallestFound[0])
                {

                }

                if (sData.sdBlocks[bi[0]].possiblePermutations.Count() < smallestFound[0])
                {

                }
            }
            Console.WriteLine($"Index: {smallestFound[2]}, count: {smallestFound[0]}, distinct: {smallestFound[2]}");

            Console.ReadLine();
        }
    }
}