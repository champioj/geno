geno
====

Geno is a golang tool to help generate generic package from an annotated package

How does it work?
-----------------

The most important component is an annotated package. It is made by adding a tag '<gen:N>' in the comment section of an interface declaration. 
You should first alias the empty interface, or any other interface you want to specialize.
eg:

    type Data interface{}  // <gen:1>
  
The new alias must be used consistently across your package.
If you need more than one specialized type, you can list them like this:

    type Key interface{}  // <gen:1>
    type Value interface{} // <gen:2>
    
// Generating

Once you have an annotated package, you can generate specialized version with the geno tool.
Geno will create a new package by replacing all instance of the annotated types with the one given in parameter.
The command

    geno -package="champioj/geno/list/intlist -types="int"
    
will create a new package named intlist with every instance of the type annotated with '<gen:1>' replaced with the int type. By convention, the name of the package is the base of the path and must be a subdirectory of the package you want to specialize.
Multiple spesialized type are separated by a comma:

    geno -package="champioj/geno/list/intlist -types="int,string"
    
It is possible to use types defined in other package. For example a list of list of int would be:

    geno -package="champioj/geno/list/intofIntlist -types="champioj/geno/list/intlist.List"
    
Import will be automatically added ... but beware of cyclic dependency!

  
