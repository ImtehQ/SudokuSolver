using System;
using System.Collections.Generic;
using System.Text;
using System.Linq;

namespace SudokuSolver
{
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

        /// <summary>
        /// return true if it can remove something
        /// </summary>
        /// <param name="index"></param>
        /// <returns></returns>
        public void CheckISValid(int index)
        {
            int removedCount = 0;

            int[] li = index.getIndexOf(Dimension.Line);
            int[] ci = index.getIndexOf(Dimension.Colum);
            int[] bi = index.getIndexOf(Dimension.Block);

            if (sdLines[li[0]].data[li[1]] > 0)
            {
//                Console.WriteLine($"=================================== Index [{index}] " +
//$"location: [{li[0]},{li[1]}]/[{ci[0]},{ci[1]}]/[{bi[0]},{bi[1]}] " +
//$"|{sdLines[li[0]].possiblePermutations.Count}" +
//$"|{sdCols[ci[0]].possiblePermutations.Count}" +
//$"|{sdBlocks[bi[0]].possiblePermutations.Count}| DEFAULT= ({sdLines[li[0]].data[li[1]]})");

                return;
            }

            bool found = false;
            List<int> alreadyCheckedNumbers = new List<int>();
            for (int i = 1; i < 10; i++)
            {
                removedCount = 0;

                
                if (i == 1)
                {
                    Console.WriteLine($"~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~");
                    //Console.WriteLine($"=================================== Index [{index}] " +
                    //    $"location: [{li[0]},{li[1]}]/[{ci[0]},{ci[1]}]/[{bi[0]},{bi[1]}] " +
                    //    $"|{sdLines[li[0]].possiblePermutations.Count}" +
                    //    $"|{sdCols[ci[0]].possiblePermutations.Count}" +
                    //    $"|{sdBlocks[bi[0]].possiblePermutations.Count}|");
                }

                if (sdLines[li[0]].data[li[1]] == i)
                {
                    Console.WriteLine($"Index [{index}] value: {i}, |L:X|C:X|B:X      ( {i} ) ");
                    Console.WriteLine($"=================================== Index [{index}] " +
                        $"location: [{li[0]},{li[1]}]/[{ci[0]},{ci[1]}]/[{bi[0]},{bi[1]}] " +
                        $"|{sdLines[li[0]].possiblePermutations.Count}" +
                        $"|{sdCols[ci[0]].possiblePermutations.Count}" +
                        $"|{sdBlocks[bi[0]].possiblePermutations.Count}| DEFAULT");
                    break;
                }

                if (!alreadyCheckedNumbers.Contains(i))
                {
                    alreadyCheckedNumbers.Add(i);
                    var matchingLineSearch = sdLines[li[0]].possiblePermutations.Where(x => x[li[1]] == i).ToList();
                    var matchingColumSearch = sdCols[ci[0]].possiblePermutations.Where(x => x[ci[1]] == i).ToList();
                    var matchingBlockSearch = sdBlocks[bi[0]].possiblePermutations.Where(x => x[bi[1]] == i).ToList();


                    if (matchingLineSearch.Count == 0 || matchingColumSearch.Count == 0 || matchingBlockSearch.Count == 0)
                    {
                        if (matchingLineSearch.Count > 0)
                            removedCount += matchingLineSearch.Count;
                        for (int m = 0; m < matchingLineSearch.Count; m++)
                        {
                            sdLines[li[0]].possiblePermutations.Remove(matchingLineSearch[m]);
                        }

                        if (matchingColumSearch.Count > 0)
                            removedCount += matchingColumSearch.Count;
                        for (int m = 0; m < matchingColumSearch.Count; m++)
                        {
                            sdCols[ci[0]].possiblePermutations.Remove(matchingColumSearch[m]);
                        }

                        if (matchingBlockSearch.Count > 0)
                            removedCount += matchingBlockSearch.Count;
                        for (int m = 0; m < matchingBlockSearch.Count; m++)
                        {
                            sdBlocks[bi[0]].possiblePermutations.Remove(matchingBlockSearch[m]);
                        }
                        Console.WriteLine($"Index [{index}] value: {i}, |L:{matchingLineSearch.Count}|C:{matchingColumSearch.Count}|B:{matchingBlockSearch.Count}      X ==== [ REMOVED {removedCount} ] ");
                    }
                    else
                    {
                        Console.WriteLine($"Index [{index}] value: {i}, |L:{matchingLineSearch.Count}|C:{matchingColumSearch.Count}|B:{matchingBlockSearch.Count}      <<----------- <{i}>");
                        found = true;
                    }

                    if (i == 9 && found == false)
                    {
                        //No solution found!
                        Console.WriteLine($"Index [{index}] , =========== >> NOT FOUND!");
                    }
                }
            }
        }

        public void CheckDoubles()
        {
            for (int y = 0; y < 9; y++)
            {
                for (int i = 0; i < 9; i++)
                {

                }
            }
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
                                if (!possiblePicks.Contains(Bval))
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
        }
    }
}
