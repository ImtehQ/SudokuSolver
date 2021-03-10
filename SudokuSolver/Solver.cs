using System;
using System.Collections.Generic;
using System.Text;
using System.Linq;
using System.Collections;

namespace SudokuSolver
{
    public class SLine
    {
        public List<int> data;
        public List<List<int>> possiblePermutations;
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

        public static List<int> GetNext(this int[] data, int startingIndex, int type)
        {
            List<int> result = new List<int>();
            if (type == 0)
            {
                for (int i = startingIndex*9; i < (startingIndex * 9) +9; i++)
                {
                    result.Add(data[i]);
                }
                return result;
            }
            if (type == 1)
            {
                for (int i = startingIndex; i < 81; i+=9)
                {
                    result.Add(data[i]);
                }
                return result;
            }
            if (type == 2)
            {
                int correctedIndex = 0;


                for (int t = 0; t < 27; t+=9)
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

                        result.Add(data[s]);
                    }
                }
                return result;
            }
            return null;
        }

        public static List<List<int>> GetPermutations(this List<int> data, List<List<int>> source)
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
                // We have one of our possible n! solutions,
                // add it to the list.
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

        static void Swap(ref int a, ref int b)
        {
            var temp = a;
            a = b;
            b = temp;
        }
    }

    class Solver
    {
        public bool Solve(int[] data)
        {
            if (!IsValidLine(data)) 
                return false;
            return true;
        }

        public bool IsValidLine(int[] data)
        {
            if (data.Distinct().Count() != 9)
                return false;
            return true;
        }


    }
}
