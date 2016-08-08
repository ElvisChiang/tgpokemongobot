#!/bin/sh

ssh -n -f $POKEDEXBOTHOST "sh -c 'cd $POKEDEXBOTPATH; nohup ./tgpokemongobot >> log.txt 2>&1 &'"
