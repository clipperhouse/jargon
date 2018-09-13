using System;

namespace Jargon
{
    public interface ILemmatizingDictionary
    {
        (string Canonical, bool Found) Lookup(string[] terms, int count);
    }
}
