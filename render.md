# The Render Pipeline

## Per Pixel

set scale to maximum
ray march:
    if hit:
        if scale is minimal return value at the point for that scale.
        reduce scale by one an ray march from hit point.
