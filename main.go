package main

import (
	"errors"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

type Vec3 struct{ x, y, z float64 }

type Camera struct {
	point, vector Vec3
	fieldOfView   float64
}

type Sphere struct {
	point, color                       Vec3
	specular, lambert, ambient, radius float64
}

type Scene struct {
	camera  Camera
	objects []Sphere
	lights  [1]Vec3
}

type Ray struct {
	point  Vec3
	vector Vec3
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	spheres := []Sphere{}
	for i := 0; i < 10; i++ {
		spheres = append(spheres, Sphere{
			Vec3{rand.Float64() * 4, rand.Float64() * 4, rand.Float64() * 4},
			Vec3{255.0, 255.0, 255.0},
			0.7, 0.3, 0.1, 0.3 + (rand.Float64() * 0.4)})
	}
	scene := Scene{
		Camera{
			Vec3{0, 1.8, 10},
			Vec3{0, 3.0, 0}, 45.0},
		spheres,
		[1]Vec3{Vec3{-30, -10, 20}}}

	render(scene, image.Point{800.0, 600.0})
}

func render(scene Scene, size image.Point) {

	im := image.NewRGBA(image.Rectangle{image.Point{0, 0}, size})

	height := float64(size.X)
	width := float64(size.Y)

	upVector := Vec3{0.0, 1.0, 0.0}
	eyeVector := unitVector(subtract(scene.camera.vector, scene.camera.point))
	vpRight := unitVector(crossProduct(eyeVector, upVector))
	vpUp := unitVector(crossProduct(vpRight, eyeVector))
	fovRadians := math.Pi * (scene.camera.fieldOfView / 2.0) / 180.0
	heightWidthRatio := height / width
	halfWidth := math.Tan(fovRadians)
	halfHeight := heightWidthRatio * halfWidth
	camerawidth := halfWidth * 2.0
	cameraheight := halfHeight * 2.0
	pixelWidth := camerawidth / (width - 1.0)
	pixelHeight := cameraheight / (height - 1.0)

	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			xComp := scale(vpRight, (float64(x)*pixelWidth)-halfWidth)
			yComp := scale(vpUp, (float64(y)*pixelHeight)-halfHeight)
			rayVector := unitVector(add3(eyeVector, xComp, yComp))
			if x == 0 && y == 0 {
				log.Printf("%v", xComp)
			}
			ray := Ray{scene.camera.point, rayVector}
			err, c := trace(ray, scene, 0)
			if err == nil {
				im.Set(x, y, color.RGBA{
					uint8(math.Min(c.x, 255)),
					uint8(math.Min(c.y, 255)),
					uint8(math.Min(c.z, 255)), 255})
			}
		}
	}

	toimg, _ := os.Create("out.png")
	defer toimg.Close()

	png.Encode(toimg, im)
}

func trace(ray Ray, scene Scene, depth int) (err error, color Vec3) {
	// if depth > 3 {
	//     return BounceErrorEror("too many bounces"), Vec3{0, 0, 0}
	// }

	distance, object := intersectScene(ray, scene)

	if distance == math.Inf(1) {
		return errors.New("miss"), Vec3{255, 255, 255}
	}

	pointAtTime := add(ray.point, scale(ray.vector, distance))

	err, col := surface(ray,
		scene,
		object,
		pointAtTime,
		sphereNormal(object.point, pointAtTime),
		depth)

	if err != nil {
		return err, Vec3{0, 0, 0}
	}

	return nil, col
}

func intersectScene(ray Ray, scene Scene) (dist float64, obj Sphere) {
	closeDist := math.Inf(1)
	closeObj := Sphere{}
	for i := 0; i < len(scene.objects); i++ {
		dist := sphereIntersection(scene.objects[i], ray)
		if dist < closeDist {
			closeDist = dist
			closeObj = scene.objects[i]
		}
	}
	return closeDist, closeObj
}

func sphereIntersection(sphere Sphere, ray Ray) float64 {
	eye_to_center := subtract(sphere.point, ray.point)
	v := dotProduct(eye_to_center, ray.vector)
	eoDot := dotProduct(eye_to_center, eye_to_center)
	discriminant := (sphere.radius * sphere.radius) - eoDot + (v * v)

	if discriminant < 0 {
		return math.Inf(1)
	} else {
		return v - math.Sqrt(discriminant)
	}
}

func surface(ray Ray, scene Scene, sphere Sphere, pointAtTime Vec3, normal Vec3, depth int) (err error, col Vec3) {

	if depth > 3 {
		return errors.New("max depth reached"), Vec3{0, 0, 0}
	}

	b := sphere.color
	c := Vec3{0, 0, 0}
	lambertAmount := 0.0

	if sphere.lambert > 0 {
		for i := 0; i < len(scene.lights); i++ {
			lightPoint := scene.lights[0]
			if isLightVisible(pointAtTime, scene, lightPoint) {

				contribution := dotProduct(unitVector(
					subtract(lightPoint, pointAtTime)), normal)

				if contribution > 0 {
					lambertAmount += contribution
				}
			}
		}
	}

	if sphere.specular > 0 {
		reflectedRay := Ray{pointAtTime, reflectThrough(ray.vector, normal)}
		err, reflectedColor := trace(reflectedRay, scene, depth+1)
		if err == nil {
			c = add(c, scale(reflectedColor, sphere.specular))
		}
	}

	return nil, add3(c, scale(b, lambertAmount*sphere.lambert),
		scale(b, sphere.ambient))
}

func isLightVisible(pt Vec3, scene Scene, light Vec3) bool {
	dist, _ := intersectScene(Ray{pt, unitVector(subtract(pt, light))}, scene)
	return dist > -0.005
}

func sphereNormal(sphere Vec3, pos Vec3) Vec3 {
	return unitVector(subtract(pos, sphere))
}

// Vector operations
func crossProduct(a Vec3, b Vec3) Vec3 {
	return Vec3{
		(a.y * b.z) - (a.z * b.y),
		(a.z * b.x) - (a.x * b.z),
		(a.x * b.y) - (a.y * b.x)}
}

func add3(a Vec3, b Vec3, c Vec3) Vec3 {
	return Vec3{
		a.x + b.x + c.x,
		a.y + b.y + c.y,
		a.z + b.z + c.z}
}

func add(a Vec3, b Vec3) Vec3 {
	return Vec3{a.x + b.x, a.y + b.y, a.z + b.z}
}

func dotProduct(a Vec3, b Vec3) float64 {
	return (a.x * b.x) + (a.y * b.y) + (a.z * b.z)
}

func unitVector(a Vec3) Vec3 {
	return scale(a, 1/length(a))
}

func scale(a Vec3, t float64) Vec3 {
	return Vec3{a.x * t, a.y * t, a.z * t}
}

func length(a Vec3) float64 {
	return math.Sqrt(dotProduct(a, a))
}

func subtract(a Vec3, b Vec3) Vec3 {
	return Vec3{a.x - b.x, a.y - b.y, a.z - b.z}
}

func reflectThrough(a Vec3, normal Vec3) Vec3 {
	d := scale(normal, dotProduct(a, normal))
	return subtract(scale(d, 2.0), a)
}
