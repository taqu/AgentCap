#include <iomanip>
#include <sstream>
#include "scene.hpp"

namespace geo {

std::string render(const Scene& scene) {
    std::ostringstream os;
    os << std::fixed << std::setprecision(2);
    for (const Shape* s : scene.by_area()) {
        Point c = s->centroid();
        os << s->name() << " area=" << s->area() << " perim=" << s->perimeter()
           << " at (" << c.x << "," << c.y << ")\n";
    }
    return os.str();
}

}  // namespace geo
