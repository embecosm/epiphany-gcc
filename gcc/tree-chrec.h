/* Chains of recurrences.
   Copyright (C) 2003, 2004 Free Software Foundation, Inc.
   Contributed by Sebastian Pop <s.pop@laposte.net>

This file is part of GCC.

GCC is free software; you can redistribute it and/or modify it under
the terms of the GNU General Public License as published by the Free
Software Foundation; either version 2, or (at your option) any later
version.

GCC is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of MERCHANTABILITY or
FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License
for more details.

You should have received a copy of the GNU General Public License
along with GCC; see the file COPYING.  If not, write to the Free
Software Foundation, 59 Temple Place - Suite 330, Boston, MA
02111-1307, USA.  */

#ifndef GCC_TREE_CHREC_H
#define GCC_TREE_CHREC_H

/* Accessors for the chains of recurrences.  */
#define CHREC_VAR(NODE)           TREE_OPERAND (NODE, 0)
#define CHREC_LEFT(NODE)          TREE_OPERAND (NODE, 1)
#define CHREC_RIGHT(NODE)         TREE_OPERAND (NODE, 2)
#define CHREC_VARIABLE(NODE)      TREE_INT_CST_LOW (CHREC_VAR (NODE))



/* The following trees are unique elements.  Thus the comparison of another 
   element to these elements should be done on the pointer to these trees, 
   and not on their value.  */

extern tree chrec_not_analyzed_yet;
extern GTY(()) tree chrec_dont_know;
extern GTY(()) tree chrec_known;

/* After having added an automatically generated element, please
   include it in the following function.  */

static inline bool
automatically_generated_chrec_p (tree chrec)
{
  return (chrec == chrec_not_analyzed_yet 
	  || chrec == chrec_dont_know
	  || chrec == chrec_known);
}

/* The tree nodes aka. CHRECs.  */

static inline bool
tree_is_chrec (tree expr)
{
  if (TREE_CODE (expr) == POLYNOMIAL_CHREC
      || automatically_generated_chrec_p (expr))
    return true;
  else
    return false;
}



/* Chrec folding functions.  */
extern tree chrec_fold_plus (tree, tree, tree);
extern tree chrec_fold_minus (tree, tree, tree);
extern tree chrec_fold_multiply (tree, tree, tree);
extern tree chrec_convert (tree, tree);
extern tree count_ev_in_wider_type (tree, tree);
extern tree chrec_type (tree);

/* Operations.  */
extern tree chrec_apply (unsigned, tree, tree);
extern tree chrec_replace_initial_condition (tree, tree);
extern tree update_initial_condition_to_origin (tree);
extern tree initial_condition (tree);
extern tree evolution_part_in_loop_num (tree, unsigned);
extern tree hide_evolution_in_other_loops_than_loop (tree, unsigned);
extern tree reset_evolution_in_loop (unsigned, tree, tree);
extern tree chrec_merge (tree, tree);

/* Observers.  */
extern bool is_multivariate_chrec (tree);
extern bool chrec_is_positive (tree, bool *);
extern bool chrec_contains_symbols (tree);
extern bool chrec_contains_symbols_defined_in_loop (tree, unsigned);
extern bool chrec_contains_undetermined (tree);
extern bool tree_contains_chrecs (tree);
extern bool evolution_function_is_affine_multivariate_p (tree);
extern bool evolution_function_is_univariate_p (tree);



/* Build a polynomial chain of recurrence.  */

static inline tree 
build_polynomial_chrec (unsigned loop_num, 
			tree left, 
			tree right)
{
  if (left == chrec_dont_know
      || right == chrec_dont_know)
    return chrec_dont_know;

  return build (POLYNOMIAL_CHREC, TREE_TYPE (left), 
		build_int_cst (NULL_TREE, loop_num), left, right);
}



/* Observers.  */

/* Determines whether CHREC is equal to zero.  */

static inline bool 
chrec_zerop (tree chrec)
{
  if (chrec == NULL_TREE)
    return false;
  
  if (TREE_CODE (chrec) == INTEGER_CST)
    return integer_zerop (chrec);
  
  return false;
}

/* Determines whether the expression CHREC is a constant.  */

static inline bool 
evolution_function_is_constant_p (tree chrec)
{
  if (chrec == NULL_TREE)
    return false;

  switch (TREE_CODE (chrec))
    {
    case INTEGER_CST:
    case REAL_CST:
      return true;
      
    default:
      return false;
    }
}

/* Determine whether the given tree is an affine evolution function or not.  */

static inline bool 
evolution_function_is_affine_p (tree chrec)
{
  if (chrec == NULL_TREE)
    return false;
  
  switch (TREE_CODE (chrec))
    {
    case POLYNOMIAL_CHREC:
      if (evolution_function_is_constant_p (CHREC_LEFT (chrec))
	  && evolution_function_is_constant_p (CHREC_RIGHT (chrec)))
	return true;
      else
	return false;
      
    default:
      return false;
    }
}

/* Determine whether the given tree is an affine or constant evolution
   function.  */

static inline bool 
evolution_function_is_affine_or_constant_p (tree chrec)
{
  return evolution_function_is_affine_p (chrec) 
    || evolution_function_is_constant_p (chrec);
}

/* Determines whether EXPR does not contains chrec expressions.  */

static inline bool
tree_does_not_contain_chrecs (tree expr)
{
  return !tree_contains_chrecs (expr);
}

/* Determines whether CHREC is a loop invariant with respect to LOOP_NUM.  
   Set the result in RES and return true when the property can be computed.  */

static inline bool
no_evolution_in_loop_p (tree chrec, unsigned loop_num, bool *res)
{
  tree scev;
  
  if (chrec == chrec_not_analyzed_yet
      || chrec == chrec_dont_know
      || chrec_contains_symbols_defined_in_loop (chrec, loop_num))
    return false;

  scev = hide_evolution_in_other_loops_than_loop (chrec, loop_num);
  *res = !tree_is_chrec (scev);
  return true;
}

#endif  /* GCC_TREE_CHREC_H  */
