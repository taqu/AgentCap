#include <iostream>
#include "scene.hpp"

int main() {
    geo::Scene scene;
    scene.add(std::make_unique<geo::Circle>("sun", geo::Point{0, 0}, 2));
    scene.add(std::make_unique<geo::Rect>("field", geo::Point{0, 0}, geo::Point{4, 3}));
    scene.add(std::make_unique<geo::Triangle>("roof", geo::Point{0, 0}, geo::Point{4, 0}, geo::Point{2, 3}));
    std::cout << geo::render(scene);
}
