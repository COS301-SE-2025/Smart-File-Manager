import pytest

def foo():
    return 3

def test_foo():
    assert foo() == 3
