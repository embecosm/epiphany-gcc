/* When -fzero-initialized-in-bss was in effect, we used to only allocate
   storage for d1.a.  */
/* { dg-do run } */
/* { dg-options "" } */
extern void abort (void);
extern void exit (int);

struct { int a; int x[]; } __attribute__((packed,aligned(sizeof (int)))) d1
  = { 0, 0 };
int d2 = 0;

int
main ()
{
  d2 = 1;
  if (sizeof (d1) != sizeof (int))
    abort ();
  if (d1.x[0] != 0)
    abort ();
  exit (0);
}
