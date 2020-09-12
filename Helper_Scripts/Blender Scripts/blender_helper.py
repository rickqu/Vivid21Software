import bpy

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