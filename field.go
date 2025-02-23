package main

type Field struct {
	x          float32
	y          float32
	row        int
	col        int
	visited    bool
	isDirTop   bool
	isDirBot   bool
	isDirRight bool
	isDirLeft  bool
}
