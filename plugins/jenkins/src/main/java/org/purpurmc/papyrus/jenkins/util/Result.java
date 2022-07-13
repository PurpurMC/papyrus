package org.purpurmc.papyrus.jenkins.util;

public class Result<T, E> {
    private final T value;
    public final E error;

    public Result(T value, E error) {
        this.value = value;
        this.error = error;
    }

    public boolean isOk() {
        return error == null;
    }

    public T getValue() {
        return value;
    }

    public E getError() {
        return error;
    }

    public static <T, E> Result<T, E> ok(T value) {
        return new Result<>(value, null);
    }

    public static <T, E> Result<T, E> error(E error) {
        return new Result<>(null, error);
    }
}
