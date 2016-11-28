/*
Package decomposition provides a hereditary base-b decomposition for the goodstein machine.

The hereditary base-b decomposition of a positive integer n is

    decompose(n) = \sum_{k=1}^{\lfloor\log_{b}(n)\rfloor} n_k \times b^{decompose(k)}

where n_k is non negative and lower than b for all k.
*/
package decomposition
