#include <cmath>
#include "geometry.hpp"

namespace geo {
namespace {
double dist(Point a, Point b) { return std::hypot(a.x - b.x, a.y - b.y); }
}

Circle::Circle(std::string name, Point c, double r) : Shape(std::move(name)), c_(c), r_(r) {}
double Circle::area() { return M_PI * r_ * r_; }
double Circle::perimeter() { return 2 * M_PI * r_; }
Point Circle::centroid() const { return c_; }

Rect::Rect(std::string name, Point min, Point max) : Shape(std::move(name)), min_(min), max_(max) {}
double Rect::area() { return (max_.x - min_.x) * (max_.y - min_.y); }
double Rect::perimeter() { return 2 * ((max_.x - min_.x) + (max_.y - min_.y)); }
Point Rect::centroid() const { return {(min_.x + max_.x) / 2, (min_.y + max_.y) / 2}; }

Triangle::Triangle(std::string name, Point a, Point b, Point c) : Shape(std::move(name)), a_(a), b_(b), c_(c) {}
double Triangle::area() {
    return std::fabs((b_.x - a_.x) * (c_.y - a_.y) - (c_.x - a_.x) * (b_.y - a_.y)) / 2;
}
double Triangle::perimeter() { return dist(a_, b_) + dist(b_, c_) + dist(c_, a_); }
Point Triangle::centroid() const { return {(a_.x + b_.x + c_.x) / 3, (a_.y + b_.y + c_.y) / 3}; }

}  // namespace geo
