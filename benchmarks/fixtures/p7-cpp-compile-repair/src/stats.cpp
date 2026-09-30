#include <algorithm>
#include <limits>
#include "scene.hpp"

namespace geo {

Summary summarize(const Scene& scene) {
    Summary s{std::numeric_limits<double>::max(), 0, 0};
    scene.for_each([&](const Shape& shape) {
        s.min_area = std::min(s.min_area, shape.area());
        s.max_area = std::max(s.max_area, shape.area());
    });
    s.mean_area = scene.size() ? scene.total_area() / scene.size() : 0;
    return s;
}

}  // namespace geo
