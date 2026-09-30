#include <algorithm>
#include "scene.hpp"

namespace geo {

void Scene::add(ShapePtr s) {
    index_[s->name()] = s.get();
    shapes_.push_back(std::move(s));
}

double Scene::total_area() const {
    double sum = 0;
    for (const auto& s : shapes_) sum += s->area();
    return sum;
}

std::vector<const Shape*> Scene::by_area() const {
    std::vector<const Shape*> out;
    for (const auto& s : shapes_) out.push_back(s.get());
    std::sort(out.begin(), out.end(), [](const Shape* a, const Shape* b) { return a->area() > b->area(); });
    return out;
}

const Shape* Scene::find(const std::string& name) const {
    auto it = index_.find(name);
    return it == index_.end() ? nullptr : it->second;
}

void Scene::for_each(const std::function<void(const Shape&)>& fn) const {
    for (const auto& s : shapes_) fn(*s);
}

}  // namespace geo
