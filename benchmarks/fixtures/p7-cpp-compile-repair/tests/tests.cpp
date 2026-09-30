#include <cmath>
#include <cstdio>
#include <memory>
#include "scene.hpp"

static int failures = 0;
#define CHECK(name, cond) do { if (cond) std::printf("ok   %s\n", name); else { failures++; std::printf("FAIL %s\n", name); } } while (0)

int main() {
    using namespace geo;
    Scene scene;
    scene.add(std::make_unique<Circle>("c", Point{0, 0}, 1));
    scene.add(std::make_unique<Rect>("r", Point{0, 0}, Point{2, 3}));
    scene.add(std::make_unique<Triangle>("t", Point{0, 0}, Point{4, 0}, Point{0, 3}));
    CHECK("circle area", std::fabs(scene.find("c")->area() - M_PI) < 1e-9);
    CHECK("rect area", scene.find("r")->area() == 6);
    CHECK("triangle area", scene.find("t")->area() == 6);
    CHECK("triangle perimeter", scene.find("t")->perimeter() == 12);
    CHECK("total area", std::fabs(scene.total_area() - (12 + M_PI)) < 1e-9);
    auto order = scene.by_area();
    CHECK("by_area size", order.size() == 3);
    CHECK("by_area last is circle", order[2]->name() == "c");
    Summary s = summarize(scene);
    CHECK("summary max", s.max_area == 6);
    CHECK("summary min", std::fabs(s.min_area - M_PI) < 1e-9);
    CHECK("render mentions roof-less scene", render(scene).find("t area=6.00") != std::string::npos);
    std::printf("%d failures\n", failures);
    return failures ? 1 : 0;
}
