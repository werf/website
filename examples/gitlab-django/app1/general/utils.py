import random

"""
This function return random ASCII-string
"""
def generate_random_string():
    out_str = ''
    for i in range(0, 16):
        a = random.randint(65, 90)
        out_str += chr(a)
    return(out_str)
