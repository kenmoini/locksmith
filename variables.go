package main

var locksmithVersion string = "0.0.1"
var readConfig *Config

// BUFFERSIZE is for copying files
var BUFFERSIZE int64 = 4096 // 4096 bits = default page size on OSX

const serverUA = "Locksmith/0.0.1"
