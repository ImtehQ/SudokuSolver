
using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Web;

namespace SudokuSolver.Logics
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

        static SDataSet sData;
        static int foundTotaal = 0;
        public static void Call()
        {
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




            int attamps = 0;
            while (!sData.isSolved)
            {
                sData.isSolved = sData.IsSolved();

                foundTotaal = 0;
                attamps++;
                int i = 0;
                if (showDebug)
                    Console.WriteLine("Checking solution attempt nr. " + attamps);


                if (attamps < 10)
                {
                    for (int x = 0; x < 9; x++)
                    {
                        for (int y = 0; y < 9; y++)
                        {
                            int[] newValue;
                            if (attamps > 5)
                                newValue = sData.CalculateNextCombination(i, true);
                            else
                                newValue = sData.CalculateNextCombination(i);
                            if (input[i] == inputSolved[i] && showDebug)
                                Console.Write($"|({inputSolved[i]}/{newValue[0]})| ");
                            else
                            {
                                if (showDebug)
                                    Console.Write($"| {inputSolved[i]}/{newValue[0]} | ");
                            }
                            foundTotaal += newValue[1];
                            i++;
                        }
                        if (showDebug)
                            Console.Write($"<{sData.sdLines[x].possiblePermutations.Count()}>");
                        if (showDebug)
                            Console.WriteLine();
                    }

                }
                else
                {
                    break;
                }
                if (showDebug)
                    Console.WriteLine($"-------------------New found: {foundTotaal}");

            }
            if (showDebug)
                Console.WriteLine("Calculating next combination...");




            Console.WriteLine("------------------------------");
            Console.WriteLine("Is the sudoku solved?");
            string solvedStr = sData.isSolved ? "yes!" : "no :(";
            Console.WriteLine($"Uuumm... : {solvedStr}");

            for (int y = 0; y < 9; y++)
            {
                for (int x = 0; x < 9; x++)
                {
                    Console.Write($"|{sData.sdLines[y].data[x]}|");
                }
                Console.WriteLine();
            }

            Console.ReadLine();
        }
    }

    public class SDataSet
    {
        public List<List<int>> Permutations;

        public SData[] sdLines = new SData[9];
        public SData[] sdCols = new SData[9];
        public SData[] sdBlocks = new SData[9];

        public bool isSolved = false;

        int[] IndexOf = new int[2];
        public void ReadSudokuData(int[] sudoku, Dimension type)
        {
            for (int i = 0; i < 81; i++)
            {
                IndexOf = i.getIndexOf(type);
                if (type == Dimension.Line)
                {
                    sdLines[IndexOf[0]].data[IndexOf[1]] = sudoku[i];
                }
                if (type == Dimension.Colum)
                {
                    sdCols[IndexOf[0]].data[IndexOf[1]] = sudoku[i];
                }
                if (type == Dimension.Block)
                {
                    sdBlocks[IndexOf[0]].data[IndexOf[1]] = sudoku[i];
                }
            }
        }



        public void Reset()
        {
            for (int i = 0; i < 9; i++)
            {

                sdLines[i] = new SData { data = new int[9] };
                sdCols[i] = new SData { data = new int[9] };
                sdBlocks[i] = new SData { data = new int[9] };

                for (int d = 0; d < 9; d++)
                {
                    sdLines[i].data[d] = 0;
                    sdCols[i].data[d] = 0;
                    sdBlocks[i].data[d] = 0;
                }
            }
        }

        /// <summary>
        /// calculates all the possible permutations for this sudoku.
        /// </summary>
        public void GetPermutations()
        {
            for (int i = 0; i < 9; i++)
            {
                sdLines[i].possiblePermutations = sdLines[i].data.GetAllPermutations(Permutations);
            }
            for (int i = 0; i < 9; i++)
            {
                sdCols[i].possiblePermutations = sdCols[i].data.GetAllPermutations(Permutations);
            }
            for (int i = 0; i < 9; i++)
            {
                sdBlocks[i].possiblePermutations = sdBlocks[i].data.GetAllPermutations(Permutations);
            }
        }



        public bool isValid()
        {
            for (int i = 0; i < 9; i++)
            {
                if (sdLines[i].isValid())
                    return false;
            }
            for (int i = 0; i < 9; i++)
            {
                if (sdCols[i].isValid())
                    return false;
            }
            for (int i = 0; i < 9; i++)
            {
                if (sdBlocks[i].isValid())
                    return false;
            }
            return true;
        }

        public bool IsSolved()
        {
            if (sdLines.Any(x => x.data.Contains(0)))
                return false;
            if (sdCols.Any(x => x.data.Contains(0)))
                return false;
            if (sdBlocks.Any(x => x.data.Contains(0)))
                return false;
            return true;
        }
        public int FindIndexOfNextEmpty(SData[] data)
        {
            for (int i = 0; i < 9; i++)
            {
                int result = IndexOfNextEmpty(data[i].data);
                if (result >= 0)
                    return result;
            }
            return -1;
        }
        public int IndexOfNextEmpty(int[] data)
        {
            return Array.IndexOf(data, 0);
        }

        int totaalFound = 0;
        bool found = false;
        List<int> possiblePicks = new List<int>();
        int[] indexes, Lindexes, Cindexes, Bindexes;
        int Lval, Cval, Bval;

        public int[] CalculateNextCombination(int rawIndex, bool randomIndex = false)
        {
            totaalFound = 0;
            found = false;
            possiblePicks.Clear();


            Lindexes = rawIndex.getIndexOf(Dimension.Line);
            if (sdLines[Lindexes[0]].data[Lindexes[1]] > 0)
                return new int[] { sdLines[Lindexes[0]].data[Lindexes[1]], totaalFound };

            for (int l = 0; l < sdLines[Lindexes[0]].possiblePermutations.Count; l++)
            {
                found = false;

                //Get Value of sdlines
                Lval = sdLines[Lindexes[0]].possiblePermutations[l][Lindexes[1]];

                //Get index of colum
                Cindexes = rawIndex.getIndexOf(Dimension.Colum);
                for (int c = 0; c < sdCols[Cindexes[0]].possiblePermutations.Count; c++)
                {
                    //get Value of sdColum
                    Cval = sdCols[Cindexes[0]].possiblePermutations[c][Cindexes[1]];
                    if (Lval == Cval) //If its the same, go next
                    {
                        //Get index of block
                        Bindexes = rawIndex.getIndexOf(Dimension.Block);

                        for (int b = 0; b < sdBlocks[Bindexes[0]].possiblePermutations.Count; b++)
                        {
                            //Get Value of block
                            Bval = sdBlocks[Bindexes[0]].possiblePermutations[b][Bindexes[1]];

                            if (Cval == Bval) //If its the same, add.
                            {
                                found = true;
                                if(!possiblePicks.Contains(Bval))
                                    possiblePicks.Add(Bval);
                                break;
                            }
                        }
                        if (found)
                            break;
                    }
                }
            }
            

            if (possiblePicks.Count == 0)
                return new int[] { 0, totaalFound };
            if (possiblePicks.Count > 1 && randomIndex)
            {
                Random r = new Random();

                indexes = rawIndex.getIndexOf(Dimension.Line);
                sdLines[indexes[0]].data[indexes[1]] = possiblePicks[0];

                indexes = rawIndex.getIndexOf(Dimension.Colum);
                sdCols[indexes[0]].data[indexes[1]] = possiblePicks[0];

                indexes = rawIndex.getIndexOf(Dimension.Block);
                sdBlocks[indexes[0]].data[indexes[1]] = possiblePicks[0];

                totaalFound++;
                GetPermutations();

                return new int[] { possiblePicks[r.Next(0, possiblePicks.Count)], totaalFound };

            }
            if (possiblePicks.Count > 1 && randomIndex == false)
            {
                return new int[] { 0, totaalFound };
            }
            

            indexes = rawIndex.getIndexOf(Dimension.Line);
            sdLines[indexes[0]].data[indexes[1]] = possiblePicks[0];

            indexes = rawIndex.getIndexOf(Dimension.Colum);
            sdCols[indexes[0]].data[indexes[1]] = possiblePicks[0];

            indexes = rawIndex.getIndexOf(Dimension.Block);
            sdBlocks[indexes[0]].data[indexes[1]] = possiblePicks[0];

            totaalFound++;
            GetPermutations();

            return new int[] { possiblePicks[0], totaalFound };
            //var result = possiblePicks.GroupBy(x => x).Select(x => new { x.Key, count = x.Count() }).OrderByDescending(x => x.count).ToList();
            //if(result != null && result.Count > 0)
            //{


            //    int finalResult = result[0].Key;
            //    indexes = rawIndex.getIndexOf(Dimension.Line);
            //    sdLines[indexes[0]].data[indexes[1]] = finalResult;

            //    indexes = rawIndex.getIndexOf(Dimension.Colum);
            //    sdCols[indexes[0]].data[indexes[1]] = finalResult;

            //    indexes = rawIndex.getIndexOf(Dimension.Block);
            //    sdBlocks[indexes[0]].data[indexes[1]] = finalResult;

            //    GetPermutations();
            //    return finalResult;
            //}
            //else
            //{
            //    //Ehhhh
            //    return 0;
            //}

        }
    }

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


    public static class SolverHelper
    {
        public static int[] getIndexOf(this int index, Dimension type)
        {
            int indexOfX, indexOfY;
            if (type == Dimension.Line)
            {
                indexOfX = index / 9;
                indexOfY = index % 9;
                return new int[] { indexOfX, indexOfY };
            }
            if (type == Dimension.Colum)
            {
                indexOfX = index / 9;
                indexOfY = index % 9;
                return new int[] { indexOfY, indexOfX };
            }
            if (type == Dimension.Block)
            {
                indexOfX = index / 9;
                indexOfY = index % 9;

                int x1 = indexOfX / 3;
                int y1 = indexOfY / 3;

                int x2 = indexOfX % 3;
                int y2 = indexOfY % 3;

                return new int[] {
                    3 * x1 + y1,
                    3 * x2 + y2 };
            }
            return null;
        }

        public static string IntToString(this List<int> data)
        {
            string returnValue = "";
            for (int i = 0; i < data.Count; i++)
            {
                returnValue += data[i].ToString() + " ";
            }
            return returnValue;
        }

        public static List<List<int>> GetAllPermutations(this int[] data, List<List<int>> source)
        {
            List<List<int>> result = null;
            for (int i = 0; i < 9; i++)
            {
                if (data[i] == 0)
                    continue;
                if (result == null)
                    result = source.Where(x => x[i] == data[i]).ToList();
                else
                    result = result.Where(x => x[i] == data[i]).ToList();
            }
            if (result == null)
                return source;
            return result;
        }

        public static List<List<int>> Permute(int[] nums)
        {
            var list = new List<List<int>>();
            return DoPermute(nums, 0, nums.Length - 1, list);
        }

        static List<List<int>> DoPermute(int[] nums, int start, int end, List<List<int>> list)
        {
            if (start == end)
            {
                list.Add(new List<int>(nums));
            }
            else
            {
                for (var i = start; i <= end; i++)
                {
                    Swap(ref nums[start], ref nums[i]);
                    DoPermute(nums, start + 1, end, list);
                    Swap(ref nums[start], ref nums[i]);
                }
            }

            return list;
        }

        static int swapInt = 0;
        static void Swap(ref int a, ref int b)
        {
            swapInt = a;
            a = b;
            b = swapInt;
        }
    }
}