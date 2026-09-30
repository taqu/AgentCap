CXX ?= g++
CXXFLAGS ?= -std=c++17 -Wall -Wextra -O1 -Iinclude
SRC = src/shapes.cpp src/scene.cpp src/render.cpp src/stats.cpp
OBJ = $(SRC:.cpp=.o)

all: scene_demo

%.o: %.cpp include/geometry.hpp include/scene.hpp
	$(CXX) $(CXXFLAGS) -c $< -o $@

scene_demo: $(OBJ) src/main.cpp
	$(CXX) $(CXXFLAGS) $(OBJ) src/main.cpp -o $@

run_tests: $(OBJ) tests/tests.cpp
	$(CXX) $(CXXFLAGS) $(OBJ) tests/tests.cpp -o $@

test: run_tests
	./run_tests

clean:
	rm -f src/*.o scene_demo run_tests

.PHONY: all test clean
