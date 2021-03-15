
using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Web;

namespace SudokuSolver.Logics
{
    public enum SolverType
    {
        SolveLogical, 
        SolveGuessing, 
        SolveEmpty
    }

    public static class Solver
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

        static SDataSet sData;

        public static void Call()
        {
            if (sData == null) sData = new SDataSet();
            sData.Reset();
            Random r = new Random();

            Console.WriteLine("resetting solver...");
            
            

            sData.Permutations = SolverHelper.Permute(new int[] { 1, 2, 3, 4, 5, 6, 7, 8, 9 });

            Console.WriteLine($"Possible line solutions: {sData.Permutations.Count}");
            Console.WriteLine($"Possible line length: {sData.Permutations[0].Count}");
            Console.WriteLine($"Check count at {1}: {sData.Permutations.Where(x => x[0] == 1).Count()}");

            Console.WriteLine("Reading sudoku input...");
            sData.ReadSudokuData(input);
            Console.WriteLine("Converted to Sudoku data set model.");
            Console.WriteLine("------------------------------");
            Console.WriteLine("Calculating possible permutations...");
            sData.GetPermutations();
            Console.WriteLine("Permutations calculated.");
            Console.WriteLine("------------------------------");
            //first find next empty field

            var rl = sData.getIndexOf(80, 0);
            var rc = sData.getIndexOf(80, 1);
            var rb = sData.getIndexOf(80, 2);

            while (!sData.isSolved)
            {
                int nextEmptyField = sData.FindIndexOfNextEmpty(sData.sdLines);
                Console.WriteLine("New field found, solving...");
                sData.CalculateNextCombination(nextEmptyField);
            }

            Console.WriteLine("Calculating next combination...");
               



            Console.WriteLine("------------------------------");
            Console.WriteLine("Is the sudoku solved?");
            string solvedStr = sData.isValid() ? "yes!" : "no :(";
            Console.WriteLine($"Uuumm... : {solvedStr}");
            Console.ReadLine();
        }
    }

    public class SDataSet
    {
        public List<List<int>> Permutations;

        public SData[] sdLines = new SData[9];
        public SData[] sdCols = new SData[9];
        public SData[] sdBlocks = new SData[9];

        int indexOfX, indexOfY;


        public bool isSolved = false;

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

        public void ReadSudokuData(int[] sudoku)
        {
            if(sdLines[0] == null)
            {
                sdLines = new SData[9];
                sdCols = new SData[9];
                sdBlocks = new SData[9];
    }
            for (int i = 0; i < 9; i++)
            {
                sdLines[i].ReadSudokuData(sudoku, i, 0);
            }
            for (int i = 0; i < 9; i++)
            {
                sdCols[i].ReadSudokuData(sudoku, i, 1);
            }
            for (int i = 0; i < 9; i++)
            {
                sdBlocks[i].ReadSudokuData(sudoku, i, 2);
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

        public bool IsSolved(int[] data)
        {
            return !data.Any(x => x == 0);
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

        public void CalculateNextCombination(int indexOfNextField)
        {
            List<int> possiblePicks = new List<int>();

            for (int i = 0; i < 9; i++)
            {
                possiblePicks.Clear();

                if (IsSolved(sdLines[i].data) == false)
                {
                    for (int d = 0; d < sdLines[i].data.Length; d++)
                    {
                        possiblePicks.Add(sdLines[i].data[d]);
                    }
                }

                for (int d = 0; d < sdCols[i].data.Length; d++)
                {
                    possiblePicks.Add(sdCols[i].data[d]);
                }
                for (int d = 0; d < sdBlocks[i].data.Length; d++) //needs a index lookup table :(
                {
                    possiblePicks.Add(sdBlocks[i].data[d]);
                }

                possiblePicks = possiblePicks.Distinct().ToList();

                if (possiblePicks.Count > 0)
                { }
            }
        }

        public int[] getIndexOf(int index, int type)
        {
            if (type == 0)
            {
                indexOfX = index / 9;
                indexOfY = index % 9;
                return new int[] { indexOfX, indexOfY };
            }
            if (type == 1)
            {
                indexOfX = index / 9;
                indexOfY = index % 9;
                return new int[] { indexOfY, indexOfX };
            }
            if(type == 2)
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
    }

    public class SData
    {
        public int[] data = new int[9];

        public List<List<int>> possiblePermutations;
        

        /// <summary>
        /// Checks is the data is valid as a sudoku line
        /// </summary>
        /// <returns></returns>
        public bool isValid()
        {
            if (!data.Contains('0'))
                return (data.Distinct().Count() == 9);
            return false;
        }

        public void ReadSudokuData(int[] sudoku, int startingIndex, int type)
        {
            List<int> result = new List<int>();
            if (type == 0)
            {
                for (int i = startingIndex * 9; i < (startingIndex * 9) + 9; i++)
                {
                    result.Add(sudoku[i]);
                }
            }
            if (type == 1)
            {
                for (int i = startingIndex; i < 81; i += 9)
                {
                    result.Add(sudoku[i]);
                }
            }
            if (type == 2)
            {
                int correctedIndex = 0;

                for (int t = 0; t < 27; t += 9)
                {
                    for (int i = 0; i < 3; i++)
                    {
                        if (startingIndex == 0)
                            correctedIndex = 0;
                        if (startingIndex == 1)
                            correctedIndex = 3;
                        if (startingIndex == 2)
                            correctedIndex = 6;
                        if (startingIndex == 3)
                            correctedIndex = 27;
                        if (startingIndex == 4)
                            correctedIndex = 30;
                        if (startingIndex == 5)
                            correctedIndex = 33;
                        if (startingIndex == 6)
                            correctedIndex = 54;
                        if (startingIndex == 7)
                            correctedIndex = 57;
                        if (startingIndex == 8)
                            correctedIndex = 60;
                        int s = correctedIndex + t + i;

                        result.Add(sudoku[s]);
                    }
                }
            }
            data = result.ToArray();
        }


    }


    public static class SolverHelper
    {
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