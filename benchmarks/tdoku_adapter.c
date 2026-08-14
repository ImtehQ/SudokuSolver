#define _GNU_SOURCE
#include "tdoku.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int main(void) {
    char *line = NULL;
    size_t capacity = 0;
    size_t puzzles = 0;
    size_t unique = 0;
    size_t guesses = 0;
    char solution[81];

    while (getline(&line, &capacity, stdin) != -1) {
        if (line[0] == '#' || strlen(line) < 81) {
            continue;
        }
        size_t count = SolveSudoku(line, 2, 0, solution, &guesses);
        puzzles++;
        if (count == 1) {
            unique++;
        } else {
            fprintf(stderr, "puzzle %zu returned %zu solution(s) with limit 2\n", puzzles, count);
            free(line);
            return 1;
        }
    }

    free(line);
    printf("{\"puzzles\":%zu,\"unique\":%zu}\n", puzzles, unique);
    return puzzles == unique ? 0 : 1;
}
