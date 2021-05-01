# newGalaxy

# Source Notes
    $ cat mk.newgal
    NEW_OBJS = NewGalaxy.o get_star.o utils.o gen_plan.o
    NewGalaxy: $(NEW_OBJS)
        cc $(NEW_OBJS) -no-pie -o ../bin/NewGalaxy
    NewGalaxy.o: NewGalaxy.c fh.h
        cc -no-pie -c NewGalaxy.c
    utils.o: utils.c fh.h
        cc -no-pie -c utils.c
    gen_plan.o: gen_plan.c fh.h
        cc -no-pie -c gen_plan.c
    get_star.o: get_star.c fh.h
        cc -no-pie -c get_star.c
    $ wc -l NewGalaxy.c utils.c gen_plan.c get_star.c
          59 get_star.c
         273 gen_plan.c
         411 NewGalaxy.c
         966 utils.c
        1709 total