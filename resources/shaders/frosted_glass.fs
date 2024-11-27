#version 330

// Input vertex attributes (from vertex shader)
in vec2 fragTexCoord;
in vec4 fragColor;

// Uniform inputs
uniform sampler2D texture0;
uniform vec4 colDiffuse;

// Custom uniforms for noise effect
uniform float noiseIntensity;  // Controls the magnitude of pixel displacement
uniform float pixelMovementPercentage;  // Percentage of pixels to be moved (0.0 to 1.0)
uniform vec2 resolution;
uniform float time;  // Optional: can be used for animated noise

// Random function
float random(vec2 st) {
    return fract(sin(dot(st.xy, vec2(12.9898, 78.233))) * 43758.5453123);
}

// Main shader function
out vec4 finalColor;

void main() {
    // Normalized pixel coordinates (from 0 to 1)
    vec2 uv = fragTexCoord;

    // Generate noise
    vec2 noiseCoord = uv;
    float noise = random(noiseCoord + time * 0.01);

    // Determine if this pixel should be moved based on probability
    if (noise < pixelMovementPercentage) {
        // Calculate random displacement vector
        vec2 displacement = vec2(
            (random(uv + vec2(1.0, 0.0)) - 0.5) * noiseIntensity,
            (random(uv + vec2(0.0, 1.0)) - 0.5) * noiseIntensity
        );

        // Apply displacement
        uv += displacement / resolution;
    }

    // Sample the texture with modified coordinates
    vec4 texColor = texture(texture0, uv);

    // Final color output
    finalColor = texColor * colDiffuse;
}