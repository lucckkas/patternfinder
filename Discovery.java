import java.math.BigInteger;
import java.util.*;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class Discovery {
    public static List<String> generarSubsecuencias(String secuenciaEntrada) {

        int largoSecuencia = secuenciaEntrada.length();

        BigInteger totalNumeros = BigInteger.ONE.shiftLeft(largoSecuencia);

        List<String> subsecuencias = new ArrayList<>();
        // recorrer todas las subsecuencias posibles
        for (BigInteger i = BigInteger.ZERO; i.compareTo(totalNumeros) < 0; i = i.add(BigInteger.ONE)) {
            StringBuilder nuevaSecuencia = new StringBuilder();

            // recorrer la secuencia de entrada
            for (int idx = 0; idx < largoSecuencia; idx++) {
                char caracter = secuenciaEntrada.charAt(idx);
                if (i.testBit(largoSecuencia - idx - 1)) {
                    nuevaSecuencia.append('x');
                } else {
                    nuevaSecuencia.append(caracter);
                }
            }
            subsecuencias.add(nuevaSecuencia.toString());

            int contadorX = 0;
            StringBuilder nuevaSecuenciaConNumeros = new StringBuilder();
            for (int j = 0; j < nuevaSecuencia.length(); j++) {
                char caracter = nuevaSecuencia.charAt(j);
                if (caracter == 'x') {
                    contadorX++;
                } else {
                    if (contadorX > 0) {
                        nuevaSecuenciaConNumeros.append(contadorX);
                        contadorX = 0;
                    }
                    nuevaSecuenciaConNumeros.append(caracter);
                }
            }
            if (contadorX > 0) {
                nuevaSecuenciaConNumeros.append(contadorX);
            }
            // resto del codigo ...
            StringBuilder nuevaSecuenciaFinal = new StringBuilder();
            for (int j = 0; j < nuevaSecuenciaConNumeros.length(); j++) {
                char caracter = nuevaSecuenciaConNumeros.charAt(j);
                if (Character.isDigit(caracter)) {
                    int valorContador = Integer.parseInt(String.valueOf(caracter));
                    if (valorContador > 0 && j > 0 && j < nuevaSecuenciaConNumeros.length() - 1 &&
                            !Character.isDigit(nuevaSecuenciaConNumeros.charAt(j - 1)) &&
                            !Character.isDigit(nuevaSecuenciaConNumeros.charAt(j + 1))) {
                        nuevaSecuenciaFinal.append(caracter);
                    }
                } else {
                    nuevaSecuenciaFinal.append(caracter);
                }
            }

            // si no tiene mayúsculas, no se agrega
            if (nuevaSecuenciaFinal.toString().matches(".*[A-Z].*")) {

                subsecuencias.add(nuevaSecuenciaFinal.toString());
            } else {

            }
        }

        return subsecuencias;
    }

    public static List<String> compararTransformaciones(List<String> lista1, List<String> lista2) {

        Set<String> set1 = new HashSet<>(lista1);
        Set<String> set2 = new HashSet<>(lista2);
        List<String> listaPatrones = new ArrayList<>();
        for (String subsequence1 : set1) {
            if (set2.contains(subsequence1)) {
                // Secuencias idénticas
                String resultado = formatIdenticalSequence(subsequence1);
                if (resultado.matches(".*[A-Z].*")) {
                    listaPatrones.add(resultado);
                }
            } else {
                for (String subsequence2 : set2) {
                    if (sequencesWithoutNumbersEqual(subsequence1, subsequence2)) {

                        // Secuencias similares (sin números)
                        String resultado = compareSequences(subsequence1, subsequence2);
                        if (resultado.matches(".*[A-Z].*")) {
                            listaPatrones.add(resultado);
                        }
                    }
                }
            }
        }

        return listaPatrones;
    }

    private static String formatIdenticalSequence(String sequence) {

        StringBuilder nuevaSubsecuencia = new StringBuilder();
        int len = sequence.length();
        for (int j = 0; j < len; j++) {
            char char1 = sequence.charAt(j);
            if (Character.isLetter(char1)) {
                if (j != 0) {
                    nuevaSubsecuencia.append("-").append(char1);
                } else {
                    nuevaSubsecuencia.append(char1);
                }
            } else if (Character.isDigit(char1)) {
                nuevaSubsecuencia.append("-x(").append(char1).append(")");
            }
        }

        return nuevaSubsecuencia.toString();
    }

    private static boolean sequencesWithoutNumbersEqual(String seq1, String seq2) {
        String seq1WithoutNumbers = seq1.replaceAll("\\d", "");
        String seq2WithoutNumbers = seq2.replaceAll("\\d", "");

        boolean iguales = seq1WithoutNumbers.equals(seq2WithoutNumbers);

        return iguales;
    }

    private static String compareSequences(String seq1, String seq2) {

        StringBuilder nuevaSubsecuencia = new StringBuilder();
        int len = Math.min(seq1.length(), seq2.length());
        int contadorX = 0;
        for (int j = 0; j < len; j++) {
            char char1 = seq1.charAt(j);
            char char2 = seq2.charAt(j);

            if (Character.isLetter(char1) && Character.isLetter(char2)) {
                if (char1 == char2) {
                    if (contadorX > 0) {
                        nuevaSubsecuencia.append("-x(").append(contadorX).append(")");
                        contadorX = 0;
                    }
                    if (j != 0) {
                        nuevaSubsecuencia.append("-").append(char1);
                    } else {
                        nuevaSubsecuencia.append(char1);
                    }

                } else {
                    contadorX++;

                }
            } else if (Character.isDigit(char1) && Character.isDigit(char2)) {
                if (char1 == char2) {
                    if (contadorX > 0) {
                        nuevaSubsecuencia.append("-x(").append(contadorX).append(")");
                        contadorX = 0;
                    }
                    nuevaSubsecuencia.append("-x(").append(char1).append(")");

                } else {
                    int num1 = Character.getNumericValue(char1);
                    int num2 = Character.getNumericValue(char2);
                    int menor = Math.min(num1, num2);
                    int mayor = Math.max(num1, num2);
                    nuevaSubsecuencia.append("-x(").append(menor).append(",").append(mayor).append(")");

                }
            } else {
                contadorX++;

            }
        }
        if (contadorX > 0) {
            nuevaSubsecuencia.append("-x(").append(contadorX).append(")");
        }

        return nuevaSubsecuencia.toString();
    }

    public static List<Integer> sumaStrings(List<String> listaStrings, int len) {

        List<Integer> resultados = new ArrayList<>();
        for (String s : listaStrings) {

            Pattern pattern = Pattern.compile("x\\((\\d+)(?:,(\\d+))?\\)");
            Matcher matcher = pattern.matcher(s);
            int suma = 0;
            while (matcher.find()) {
                int valor1 = Integer.parseInt(matcher.group(1));
                int valor2 = matcher.group(2) != null ? Integer.parseInt(matcher.group(2)) : 0;
                int maxValor = Math.max(valor1, valor2);
                suma += maxValor;
            }
            for (char c : s.toCharArray()) {
                if (Character.isLetter(c) && c != 'x') {
                    if (Character.isUpperCase(c)) {
                        suma += len * len;
                    } else {
                        suma += len;
                    }

                }
            }
            resultados.add(suma);

        }

        return resultados;
    }

    public static List<String> ordenarListaConNumeros(List<String> listaStrings, List<Integer> listaNumeros) {

        List<String> listaOrdenada = new ArrayList<>();
        List<Integer> indices = new ArrayList<>();
        for (int i = 0; i < listaStrings.size(); i++) {
            indices.add(i);
        }

        Collections.sort(indices, (i1, i2) -> {
            return Integer.compare(listaNumeros.get(i2), listaNumeros.get(i1));
        });

        for (int indice : indices) {
            listaOrdenada.add(listaStrings.get(indice));
        }

        return listaOrdenada;
    }

    private static String leerSecuenciasDesdeConsola(Scanner scanner) {
        String secuencia = new String();
        while (scanner.hasNextLine()) {
            String linea = scanner.nextLine();
            if (!linea.isEmpty()) {
                System.out.println("Secuencia ingresada: " + linea);
                secuencia = linea;
                break;
            } else {
                System.out.println("Ingrese una secuencia válida");
            }
        }
        return secuencia;
    }

    public static void main(String[] args) {
        Scanner scanner = new Scanner(System.in);
        String secuencia1;
        String secuencia2;
        if (args.length >= 2) {
            secuencia1 = args[0];
            secuencia2 = args[1];
        } else {
            System.out.println("Ingrese secuencia 1:");
            secuencia1 = leerSecuenciasDesdeConsola(scanner);
            System.out.println("Ingrese secuencia 2:");
            secuencia2 = leerSecuenciasDesdeConsola(scanner);
        }
        System.out.println("Secuencia 1: " + secuencia1);
        System.out.println("Secuencia 2: " + secuencia2);
        long startTime = System.nanoTime();
        List<String> lista1Sec = generarSubsecuencias(secuencia1);
        List<String> lista2Sec = generarSubsecuencias(secuencia2);
        List<String> listaPatrones = compararTransformaciones(lista1Sec, lista2Sec);
        System.out.println("Patrones encontrados: " + listaPatrones);
        int len = Math.max(secuencia1.length(), secuencia2.length());
        List<Integer> resultados = sumaStrings(listaPatrones, len);
        List<String> listaOrdenada = ordenarListaConNumeros(listaPatrones, resultados);
        long endTime = System.nanoTime();
        // Convertir a milisegundos
        long duration = (endTime - startTime) / 1000000;
        System.out.println("Patrón final para el par (" + secuencia1 + ", " + secuencia2
                + "): " + listaOrdenada);
        System.out.println("Duración: " + duration + " ms");
        scanner.close();
    }
}
