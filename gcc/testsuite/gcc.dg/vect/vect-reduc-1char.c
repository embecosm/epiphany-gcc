/* { dg-require-effective-target vect_int } */

#include <stdarg.h>
#include "tree-vect.h"

#define N 16
#define DIFF 242

int main1 (unsigned char x, unsigned char max_result)
{
  int i;
  unsigned char ub[N] = {0,3,6,9,12,15,18,21,24,27,30,33,36,39,42,45};
  unsigned char uc[N] = {0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15};
  unsigned char udiff = 2;
  unsigned char umax = x;
  unsigned char umin = 10;

  for (i = 0; i < N; i++) {
    udiff += (unsigned char)(ub[i] - uc[i]);
  }

  for (i = 0; i < N; i++) {
    umax = umax < uc[i] ? uc[i] : umax;
  }

  for (i = 0; i < N; i++) {
    umin = umin > uc[i] ? uc[i] : umin;
  }

  /* check results:  */
  if (udiff != DIFF)
    abort ();
  if (umax != max_result)
    abort ();
  if (umin != 0)
    abort ();

  return 0;
}

int main (void)
{ 
  check_vect ();
  
  main1 (100, 100);
  main1 (0, 15);
  return 0;
}

/* { dg-final { scan-tree-dump-times "vectorized 3 loops" 1 "vect" } } */
/* { dg-final { cleanup-tree-dump "vect" } } */
