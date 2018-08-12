class Factorial {
  public static int factorial(int x) {
    return x == 1 ? 1 : x * factorial(x - 1);
  }
  public static void main(String[] args) {
    factorial(10);
  }
}
