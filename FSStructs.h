//
//  FSStructs.h
//  HoteApplication
//
//  Created by Adel on 07/01/2021.
//  Copyright © 2021 ABTasty. All rights reserved.
//

#ifndef FSStructs_h
#define FSStructs_h

#include <stdio.h>


typedef struct {

    char *msg;
}FSContext;

typedef struct {

    int a;

}FSModifications;


typedef struct userContextBis {
    char *name;
    int var_type; // 1 = string, 2 = bool, 3 = float, 4 = int
    void *data;
} userContextBis;





#endif /* FSStructs_h */
