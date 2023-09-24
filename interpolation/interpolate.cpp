#include <NTL/ZZ_pX.h>
#include <iostream>
#include <vector>
#include <sstream>

using namespace NTL;
using namespace std;

extern "C" {
    const char* interpolate(const char* prime, int n, const char** x_vals, const char** y_vals);
}

const char* interpolate(const char* prime, int n, const char** x_vals, const char** y_vals) {
    static string result;

    ZZ p = conv<ZZ>(prime);
    ZZ_p::init(p);

    vector<ZZ_p> x(n), y(n);
    for (int i = 0; i < n; i++) {
        x[i] = conv<ZZ_p>(conv<ZZ>(x_vals[i]));
        y[i] = conv<ZZ_p>(conv<ZZ>(y_vals[i]));
    }

    ZZ_pX P;
    clear(P);

    for (int i = 0; i < n; i++) {
        ZZ_pX Li;
        Li = 1;
        for (int j = 0; j < n; j++) {
            if (i != j) {
                Li *= (ZZ_pX(1, 1) - x[j]) / (x[i] - x[j]);
            }
        }
        P += y[i] * Li;
    }

    stringstream ss;
    ss << P;
    result = ss.str();
    return result.c_str();
}
