#!/usr/bin/python

import base
import sys

class RavagePlayer(base.AbstractTrickTakingPlayer):
  def __init__(self, hand_cards):
    base.AbstractTrickTakingPlayer.__init__(self, hand_cards)
    self._taken_by_suit = [0, 0, 0, 0]

  def GetHandResultDescription(self):
    return str(self._taken_by_suit)

  def RecordTrick(self, cards, my_card_index):
    winner = base.AbstractTrickTakingPlayer.RecordTrick(
      self, cards, my_card_index)
    if winner == cards[my_card_index]:
      for card in cards:
        self._taken_by_suit[card.suit] += 1

  def _ChooseLead(self, suit):
    my_values = self.hand.GetValues(suit)
    return (my_values[0], 20 - my_values[0])

  def _ChooseFollow(self, suit, played_values, discards):
    num_played = len(played_values) + len(discards)
    highest_played = max(played_values)
    my_values = self.hand.GetValues(suit)
    highest_duck_value = self._ChooseHighestDuckValue(my_values, highest_played)
    if highest_duck_value is not None: return highest_duck_value
    if num_played == 3: return my_values[-1]
    lowest_win_value = self._ChooseLowestWinValue(my_values, highest_played)
    assert lowest_win_value is not None
    return lowest_win_value

  def _ChooseHighestDuckValue(self, my_values, highest_played):
    if my_values[0] > highest_played: return None
    i = 0
    while i + 1 < len(my_values) and my_values[i + 1] < highest_played:
      i += 1
    return my_values[i]

  def _ChooseLowestWinValue(self, my_values, highest_played):
    i = 0
    while my_values[i] < highest_played:
      i += 1
      if i == len(my_values): return None
    return my_values[i]

  def _ChooseDiscard(self, discard_suit, suit, played_values, discards):
    my_values = self.hand.GetValues(discard_suit)
    return (my_values[-1], 20 - len(my_values))


hand_cards = base.Card.ListFromString(sys.stdin.readline())
base.RunTrickTakingPlayer(RavagePlayer(hand_cards), print_history=False)
