{-# LANGUAGE OverloadedStrings #-}
module Main where

import           Data.Time

main :: IO ()

main = do
    start <- getCurrentTime
    print "press Enter to stop your stopwatch!"
    input <- getChar
    now   <- getCurrentTime
    putStrLn $ formatTime defaultTimeLocale "%h:%m:%s" (diffUTCTime now start)
