geno
====

Geno is a golang tool to help generate generic package from an annotated package.
It's goal is to experiment/reproduce the idea of generic pacakge developped in [this topic](https://groups.google.com/forum/#!searchin/golang-nuts/generic/golang-nuts/7G2CrXNhDI0/C-bAdXdbe9kJ)

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
    
Generating
----------

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

Autogeneration
--------------

It's possible to add annotation to import so a single call to geno is able to generate every specialized package you need.
Add the import as if the package were already created. Then annotate it with its specialization:

    import(
        "champioj/geno/list/intlist" // <gen:int>
        "champioj/geno/list/intofIntlist" // <gen:champioj/geno/list/intlist.List>
    )
    
You can observe that it mimic the parameters -package and -types from the geno tool.
Once it done, you can make a single call to geno:

    geno -recursive="user/foo/bar"
    
Geno will parse the package "user/foo/bar" and all it import recursively and automatically generate every specialized type it found.

A word on cyclic dependencies
---------

The main drawback which is innerent to the logic of the progam is the problem of cyclic dependencies.
It is not possible to declare a type in a package, and use a specialized package which use that type cause it would create a cyclic dependency. If you want to do it, you have to put your type in its own packge.
I would gladly receive any idea on how to mitigate or solve this problem.


