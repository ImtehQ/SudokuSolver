using System;
using System.Collections.Generic;
using System.Text;

namespace SudokuSolver
{
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
