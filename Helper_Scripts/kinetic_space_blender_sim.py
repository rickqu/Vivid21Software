import bpy
import collections
from math import sqrt
from numpy import arange

''' A small script to render and test visuals for the "Kinetic Space" idea for VIVID 2021.

    The script uses the blender python API to generate sphere primitives in 3d space. The intent of the code is to have an accessible sandbox for testing and visualising the various effects.

'''

# Setup code
point = collections.namedtuple('point',['x','y'])
colour = collections.namedtuple('colour',['hue','saturation','lightness'])

# Array code with examples
def make_array(array_gen, x_min, x_max, y_min, y_max):
    ''' Sets an x-y bounding box for points to generate. The array_gen will decide how the points are generated (should return a stream of x, y)
    '''
    points = list(point(x, y) for (x, y) in array_gen(x_min, x_max, y_min, y_max))
    return points

def rec_array(spacing_x, spacing_y):
    ''' Generates a rectangular grid in the bounds using x and y spacing
    '''
    def rec_gen(x_min, x_max, y_min, y_max):
        for x in arange(x_min, x_max, spacing_x):
            for y in arange(y_min, y_max, spacing_y):
                yield x, y
    return rec_gen
    
def tri_array(spacing):
    ''' Generates a grid using equilateral triangle spacing
    '''
    def tri_gen(x_min, x_max, y_min, y_max):
        for ii, y in enumerate(arange(y_min, y_max, spacing*sqrt(3)/2)):
            for x in arange(x + spacing*(ii % 2)/2, x_max, spacing):
                yield x, y
    return tri_gen

# Blender functions
def clear_objects():
    ''' Unlinks and clears all objects in the current context/collection.
        Tested and working.
    '''
    for obj in bpy.context.collection.objects:
        bpy.context.collection.objects.unlink(obj)
        bpy.data.objects.remove(obj)

def clear_materials():
    ''' Deletes all matereials in the current dataset.
        Tested and working.
    '''
    for _, mat in bpy.data.materials.items():
        bpy.data.materials.remove(mat)

def make_sphere(rad , loc):
    ''' Makes a sphere with radius 'rad' and location 'loc' in (x, y, z).
        Returns the newly created sphere object.
        Tested and working.
    '''
    old_keys = set(bpy.context.collection.objects.keys())
    bpy.ops.mesh.primitive_ico_sphere_add(radius = rad, location = loc)
    new_keys = set(bpy.context.collection.objects.keys())
    (newest_key,) = new_keys - old_keys
    return bpy.context.collection.objects[newest_key]

def make_material(name):
    ''' Makes a material with name 'name'. If duplicate name, raises KeyError
        Returns the newly creater material object.
        Tested and working.
    '''
    if name in bpy.data.materials.keys():
        raise KeyError('Duplicate Key')
    else:
        bpy.data.materials.new(name)
    return bpy.data.materials[name]

# Key functions to play with for simulation
def z_fun(x, y):
    ''' Returns a custom height function for each point
    '''
    return x**2 + y**2

def col_fun(x, y):
    ''' Returns a colour value for each point using the HLS scale
    '''
    result_colour = colour
    result_colour.hue = (x + y) % 360
    result_colour.saturation = 1
    result_colour.lightness = 1
    return result_colour

def gen_lights(points, l_radius, z_fun, col_fun):
    ''' Calls the blender API to produce light spheres with radius l_radius at x, y points
        The points will have height z = z_fun(x, y)
        The points will have colour c = col_fun(*params)
    '''
    for p in points:
        z = z_fun(p.x, p.y)
        c = col_fun(p.x, p.y)
        
        bpy.ops.mesh.primitive_ico_sphere_add(
            location = (p.x, p.y, z),
            size = l_radius
        )
        
        # TODO: add functionality to change colour and luminosity

# Example code usage
x_min = y_min = 0
x_max = y_max = 1000
points = make_array(rec_array(200, 200), x_min, x_max, y_min, y_max)
gen_lights(points, l_radius = 120, z_fun = z_fun, col_fun = col_fun)
