#include <iostream>
#include <cuda_runtime.h>
#include <device_launch_parameters.h>

using namespace std;
const int M = 8;
const int N = 8;

__global__ void matrix_add(int *A, int *B, int *C, int width) {
    int row = blockIdx.y * blockDim.y + threadIdx.y;
    int col = blockIdx.x * blockDim.x + threadIdx.x;
    int index = row * width + col;
    if (row < M && col < N) {
        C[index] = A[index] + B[index];
    }
}

int main() {
    const int num_elements = M * N;
    const int nbytes = num_elements * sizeof(int);
    
    int *host_A = (int *)malloc(nbytes);
    int *host_B = (int *)malloc(nbytes);
    int *host_C = (int *)malloc(nbytes);

    for (int i = 0; i < num_elements; i++) {
        host_A[i] = i;
        host_B[i] = i;
    }

    int *dev_A, *dev_B, *dev_C;
    cudaMalloc((void **)&dev_A, nbytes);
    cudaMalloc((void **)&dev_B, nbytes);
    cudaMalloc((void **)&dev_C, nbytes);

    cudaMemcpy(dev_A, host_A, nbytes, cudaMemcpyHostToDevice);
    cudaMemcpy(dev_B, host_B, nbytes, cudaMemcpyHostToDevice);
    cudaMemset(dev_C, 0, nbytes);

    dim3 threadsPerBlock(8, 8);
    dim3 numBlocks(N / threadsPerBlock.x, M / threadsPerBlock.y);
    
    matrix_add<<<numBlocks, threadsPerBlock>>>(dev_A, dev_B, dev_C, N);

    cudaDeviceSynchronize();
    cudaMemcpy(host_C, dev_C, nbytes, cudaMemcpyDeviceToHost);
    for (int i = 0; i < M; i++) {
        for (int j = 0; j < N; j++) {
            cout << host_C[i * N + j] << " ";
        }
        cout << endl;
    }

    free(host_A);
    free(host_B);
    free(host_C);
    cudaFree(dev_A);
    cudaFree(dev_B);
    cudaFree(dev_C);
    return 0;
}
