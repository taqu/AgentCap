#pragma once
#include <functional>
#include <map>
#include <string>
#include <vector>
#include "geometry.hpp"

namespace geo {

class Scene {
public:
    void add(ShapePtr s);
    std::size_t size() const { return shapes_.size(); }
    double total_area() const;
    // Returns shapes sorted by descending area.
    std::vector<const Shape*> by_area() const;
    const Shape* find(const std::string& name) const;
    void for_each(const std::function<void(const Shape&)>& fn) const;

private:
    std::vector<ShapePtr> shapes_;
    std::map<std::string, const Shape*> index_;
};

std::string render(const Scene& scene);

struct Summary {
    double min_area, max_area, mean_area;
};
Summary summarize(const Scene& scene);

}  // namespace geo
