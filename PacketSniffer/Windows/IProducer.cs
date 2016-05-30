using System;
namespace Guardian
{
    interface IProducer<T>
    {
        bool Produce(T item);
    }
}
