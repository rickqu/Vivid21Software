# Browser Renderer

## Package Intent
A browser-based renderer for lighting effects in 2d

## Pseudocode

```python
def main():

    # Setup code
    input_device = some_device.start()
    boundaries = some_data_input()
    socket = make_web_socket()
    virtual_objects = initialise_virtual_objects()

    while True:
        real_objects = input_device.get_data()
        virtual_objects = update_physics_objects(real_objects, virtual_objects)
        virtual_objects = vary_virtual_objects(virtual_objects, boundaries)

        # First maps the colour positions of real objects
        mapped = map_objects(render_grid, real_objects, real_colour_function())

        # Next overrides the map if necessary with virtual objects
        mapped = map_objects(render_grid, virtual_objects, virtual_colour_function, mapped

        with socket.write(some params) as s:
            s.write(wrapped(mapped))

        with interrupt.read(some params) as inp:
            if inp == 'Stop':
                break

# Helper functions below
def update_physics_objects(real_objects, virtual_objects):
    ''' Does some physics simulation'''
    return new_virtual_objects

def vary_virtual_objects(virtual_objects, boundaries):
    new_virtual_objects = virtual_objects.copy()
    for obj in virtual_objects:
        if obj is not in boundaries and random.choice(True, False):
            new_virtual_objects.delete(obj)
            new_virtual_objects.add(make_new_physics_object(boundaries))
    return new_virtual_objects

def render_grid(spacing, x_min, x_max, y_min, y_max):
    ''' Generates a square grid within a bounding box'''
    points = []
    for x in (x_min to x_max, diff spacing):
        for y in (y_min to y_max, diff spacing):
            points.append(point(x, y))
    return points

def map_objects(render_grid, objects, colour_function, mapped = None):
    ''' Assigns colours to points in the render_grid based on what object and function'''
    if mapped is None:
        mapped = {}
    
    for point in render_grid:
        if point in objects.boundaries:
            mapped[point] = colour_function()
    return mapped

def real_colour_function()
    ''' TODO: make function mapping real objects to some colour scheme'''

def virtual_colour_function()
    ''' TODO: make function mapping virtual objects to some other colour scheme'''

```