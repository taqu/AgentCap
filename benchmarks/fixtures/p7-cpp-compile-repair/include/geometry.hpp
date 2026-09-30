#pragma once
#include <memory>
#include <string>

namespace geo {

struct Point {
    double x = 0, y = 0;
};

// Shape is the polymorphic base of every drawable primitive. All queries
// are const: shapes are immutable once placed in a Scene.
class Shape {
public:
    explicit Shape(std::string name) : name_(std::move(name)) {}
    virtual ~Shape() = default;
    virtual double area() const = 0;
    virtual double perimeter() const = 0;
    virtual Point centroid() const = 0;
    const std::string& name() const { return name_; }

private:
    std::string name_;
};

class Circle : public Shape {
public:
    Circle(std::string name, Point c, double r);
    double area() override;
    double perimeter() override;
    Point centroid() const override;

private:
    Point c_;
    double r_;
};

class Rect : public Shape {
public:
    Rect(std::string name, Point min, Point max);
    double area() override;
    double perimeter() override;
    Point centroid() const override;

private:
    Point min_, max_;
};

class Triangle : public Shape {
public:
    Triangle(std::string name, Point a, Point b, Point c);
    double area() override;
    double perimeter() override;
    Point centroid() const override;

private:
    Point a_, b_, c_;
};

using ShapePtr = std::unique_ptr<Shape>;

}  // namespace geo
