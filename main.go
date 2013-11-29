package main

import (
    "fmt"
    "math"
    "image"
    _ "image/png")

type Vec3 struct { x, y, z float64 }

type Camera struct {
    point, vector Vec3
    fieldOfView float64
}

type Sphere struct {
    point, color Vec3
    specular, lambert, ambient, radius float64
}

type Scene struct {
    camera Camera
    objects []Sphere
}

type Ray struct {
    point Vec3
    vector Vec3
}

func main() {
    im := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{800, 600}})
    a := Vec3{1, 1, 1}
    b := Vec3{1, 2, 1}
    c := dotProduct(a, b)
    fmt.Println("a:", c)
}

func render(scene Scene, size image.Point) {
    upVector := Vec3{1, 0, 0}
    eyeVector := unitVector(subtract(scene.camera.vector, scene.camera.point))
    vpRight := unitVector(crossProduct(eyeVector, upVector))
    vpUp := unitVector(crossProduct(vpRight, eyeVector))
    fovRadians := math.Pi * (scene.camera.fieldOfView / 2) * 180
    heightWidthRatio := float64(size.Y) / float64(size.X)
    halfWidth := math.Tan(fovRadians)
    halfHeight := heightWidthRatio * halfWidth
    camerawidth := halfWidth * 2
    cameraheight := halfHeight * 2
    pixelWidth := camerawidth / (float64(size.X) - 1)
    pixelHeight := cameraheight / (float64(size.Y) - 1)

    for x := 0; x < size.X; x++ {
        for y := 0; y < size.Y; y++ {
            xComp := scale(vpRight, (float64(x) * pixelWidth - halfWidth))
            yComp := scale(vpUp, (float64(y) * pixelHeight) - halfHeight)
            ray := Ray(camera.point, unitVector(add3(eyeVector, xComp, yComp)))
            color := trace(ray, scene, 0)
        }
    }
}

func trace(ray Vec3, scene Scene, depth int) (err error, color Vec3) {
    // if depth > 3 {
    //     return BounceErrorEror("too many bounces"), Vec3{0, 0, 0}
    // }

    distance, object := intersectScene(ray, scene)

    if distance == math.Inf(1) {
        return nil, Vec3{255, 255, 255}
    }

    pointAtTime := add(ray.point, scale(ray.vector, distance))
}

func intersectScene(ray Vec3, scene Scene) (dist float64, obj Sphere) {
    closeDist := math.Inf(1)
    closeObj := nil
    for obj := range scene.objects {
        dist := sphereIntersection(obj, ray)
        if dist < closeDist {
            closeDist := dist
            closeObj := obj
        }
    }
    return closeDist, closeObj
}

func sphereIntersection(sphere Sphere, ray Vec3) float64 {
    eye_to_center := subtract(sphere.point, ray.point)
    v := dotProduct(eye_to_center, ray.vector)
    eoDot := dotProduct(eye_to_center, eye_to_center)
    discriminant := (sphere.radius * sphere.radius) - eoDot + (v * v)

    if discriminant < 0 {
        return nil
    } else {
        return v - math.Sqrt(discriminant)
    }
}

func isLightVisible (pt Vec3) bool {
    return false
}

func sphereNormal(sphere Vec3, pos Vec3) Vec3 {
    return unitVector(subtract(pos, sphere))
}

// Vector operations
func crossProduct (a Vec3, b Vec3) Vec3 {
    return Vec3{
        (a.y * b.z) - (a.z * b.y),
        (a.z * b.x) - (a.x * b.z),
        (a.x * b.y) - (a.y * b.x) }
}

func add3 (a Vec3, b Vec3, c Vec3) Vec3 {
    return Vec3{
        a.x + b.x + c.x,
        a.y + b.y + c.y,
        a.z + b.z + c.z }
}

func add (a Vec3, b Vec3) Vec3 {
    return Vec3{
        a.x + b.x,
        a.y + b.y,
        a.z + b.z }
}

func dotProduct (a Vec3, b Vec3) float64 {
    return (a.x * b.x) + (a.y * b.y) + (a.z * b.z)
}

func unitVector (a Vec3) Vec3 {
    return scale(a, 1 / length(a))
}

func scale (a Vec3, t float64) Vec3 {
    return Vec3{a.x * t, a.y * t, a.z * t }
}

func length (a Vec3) float64 {
    return math.Sqrt(dotProduct(a, a))
}

func subtract (a Vec3, b Vec3) Vec3 {
    return Vec3{ a.x - b.x, a.y - b.y, a.z - b.z }
}
