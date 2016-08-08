#!/bin/sh

ssh $POKEDEXBOTHOST "sh -c 'tail -f $POKEDEXBOTPATH/log.txt'"
