using System.Collections.Generic;

namespace Jargon.Data
{
    public class Numbers : ILemmatizingDictionary
    {
        public static readonly Numbers Instance = new Numbers();

        private static readonly Dictionary<string, long> NumbersLookup = new Dictionary<string, long>
            {
                ["one"] = 1,
                ["two"] = 2,
                ["three"] = 3,
                ["four"] = 4,
                ["five"] = 5,
                ["six"] = 6,
                ["seven"] = 7,
                ["eight"] = 8,
                ["nine"] = 9,
                ["ten"] = 10,
                ["eleven"] = 11,
                ["twelve"] = 12,
                ["thirteen"] = 13,
                ["fourteen"] = 14,
                ["fifteen"] = 15,
                ["sixteen"] = 16,
                ["seventeen"] = 17,
                ["eighteen"] = 18,
                ["nineteen"] = 19,
                ["twenty"] = 20,
                ["twentyone"] = 21,
                ["twentytwo"] = 22,
                ["twentythree"] = 23,
                ["twentyfour"] = 24,
                ["twentyfive"] = 25,
                ["twentysix"] = 26,
                ["twentyseven"] = 27,
                ["twentyeight"] = 28,
                ["twentynine"] = 29,
                ["thirty"] = 30,
                ["thirtyone"] = 31,
                ["thirtytwo"] = 32,
                ["thirtythree"] = 33,
                ["thirtyfour"] = 34,
                ["thirtyfive"] = 35,
                ["thirtysix"] = 36,
                ["thirtyseven"] = 37,
                ["thirtyeight"] = 38,
                ["thirtynine"] = 39,
                ["forty"] = 40,
                ["fortyone"] = 41,
                ["fortytwo"] = 42,
                ["fortythree"] = 43,
                ["fortyfour"] = 44,
                ["fortyfive"] = 45,
                ["fortysix"] = 46,
                ["fortyseven"] = 47,
                ["fortyeight"] = 48,
                ["fortynine"] = 49,
                ["fifty"] = 50,
                ["fiftyone"] = 51,
                ["fiftytwo"] = 52,
                ["fiftythree"] = 53,
                ["fiftyfour"] = 54,
                ["fiftyfive"] = 55,
                ["fiftysix"] = 56,
                ["fiftyseven"] = 57,
                ["fiftyeight"] = 58,
                ["fiftynine"] = 59,
                ["sixty"] = 60,
                ["sixtyone"] = 61,
                ["sixtytwo"] = 62,
                ["sixtythree"] = 63,
                ["sixtyfour"] = 64,
                ["sixtyfive"] = 65,
                ["sixtysix"] = 66,
                ["sixtyseven"] = 67,
                ["sixtyeight"] = 68,
                ["sixtynine"] = 69,
                ["seventy"] = 70,
                ["seventyone"] = 71,
                ["seventytwo"] = 72,
                ["seventythree"] = 73,
                ["seventyfour"] = 74,
                ["seventyfive"] = 75,
                ["seventysix"] = 76,
                ["seventyseven"] = 77,
                ["seventyeight"] = 78,
                ["seventynine"] = 79,
                ["eighty"] = 80,
                ["eightyone"] = 81,
                ["eightytwo"] = 82,
                ["eightythree"] = 83,
                ["eightyfour"] = 84,
                ["eightyfive"] = 85,
                ["eightysix"] = 86,
                ["eightyseven"] = 87,
                ["eightyeight"] = 88,
                ["eightynine"] = 89,
                ["ninety"] = 90,
                ["ninetyone"] = 91,
                ["ninetytwo"] = 92,
                ["ninetythree"] = 93,
                ["ninetyfour"] = 94,
                ["ninetyfive"] = 95,
                ["ninetysix"] = 96,
                ["ninetyseven"] = 97,
                ["ninetyeight"] = 98,
                ["ninetynine"] = 99,
            };
        private static readonly Dictionary<string, long> Magnitudes = new Dictionary<string, long>
        {
            ["hundred"]     = 100,
            ["thousand"]    = 1_000,
            ["million"]     = 1_000_000,
            ["billion"]     = 1_000_000_000,
            ["trillion"]    = 1_000_000_000_000,
            ["quadrillion"] = 1_000_000_000_000_000,
            ["quintillion"] = 1_000_000_000_000_000_000
        };

        private Numbers() { }

        public (string Canonical, bool Found) Lookup(string[] terms, int count)
        {
            if(count == 0)
            {
                return ("", false);
            }

            var p = new Parser(terms, count);
            return p.Parse();
        }

        private class Parser
        {
            private string[] Tokens;
            private int Count;

            private List<long> Ints;
            private List<double> Floats;
            private int Pos;

            public Parser(string[] tokens, int count)
            {
                Tokens = tokens;
                Count = count;

                Ints = null;
                Floats = null;
                Pos = 0;
            }

            public (string Canonical, bool Found) Parse()
            {
                var p = this;
                p.Pos = 0;

                var first = p.Current();

                var success = long.TryParse(first, out var i);
                if (success)
                {
                    if(HasLeadingZero(first))
                    {
                        return ("", false);
                    }
                    else
                    {
                        if(p.Ints == null)
                        {
                            p.Ints = new List<long>();
                        }
                        p.Ints.Add(i);
                        p.Pos++;
                        return p.ParseMagnitudesInt();
                    }
                }

                success = double.TryParse(first, out var f);
                if (success)
                {
                    if(p.Floats == null)
                    {
                        p.Floats = new List<double>();
                    }
                    p.Floats.Add(f);
                    p.Pos++;
                    return p.ParseMagnitudesFloat();
                }

                success = NumbersLookup.TryGetValue(first, out var num);
                if (success)
                {
                    if(p.Ints == null)
                    {
                        p.Ints = new List<long>();
                    }
                    p.Ints.Add(num);
                    p.Pos++;
                    return p.ParseMagnitudesInt();
                }

                return ("", false);
            }

            private string Current() => Normalize(this.Tokens[this.Pos]);

            private (string Canonical, bool Found) ParseMagnitudesInt()
            {
                var p = this;

                while(p.Pos < p.Count)
                {
                    var ok = Magnitudes.TryGetValue(p.Current(), out var m);
                    if (!ok)
                    {
                        return ("", false);
                    }

                    if(p.Ints == null)
                    {
                        p.Ints = new List<long>();
                    }
                    p.Ints.Add(m);
                    p.Pos++;
                }

                var result = 1L;
                for(var i = 0; i < p.Ints.Count; i++)
                {
                    result = result * p.Ints[i];
                }

                return (result.ToString(), true);
            }

            private (string Canonical, bool Found) ParseMagnitudesFloat()
            {
                var p = this;

                while (p.Pos < p.Count)
                {
                    var ok = Magnitudes.TryGetValue(p.Current(), out var m);
                    if (!ok)
                    {
                        return ("", false);
                    }

                    if (p.Floats == null)
                    {
                        p.Floats = new List<double>();
                    }
                    p.Floats.Add(m);
                    p.Pos++;
                }

                var result = 1.0;
                for (var i = 0; i < p.Floats.Count; i++)
                {
                    result = result * p.Floats[i];
                }
                
                return (result.ToString("R"), true);
            }

            private static string Normalize(string s)
            {
                var result = s;
                result = result.Replace(",", "");
                result = result.Replace("-", "");

                if (s.Length > 0 && s[0] == '-')
                {
                    return "-" + result;
                }

                return result;
            }

            private static bool HasLeadingZero(string s)
            {
                if (s.Length == 0) return false;

                if (s[0] == '0') return true;

                if(s.Length > 1)
                {
                    if ((s[0] == '-' || s[0] == '+') && s[1] == '0') return true;
                }

                return false;
            }
        }
    }
}
